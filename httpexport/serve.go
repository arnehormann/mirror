package httpexport

import (
	"bufio"
	"fmt"
	"github.com/arnehormann/mirror"
	"net/http"
	"reflect"
)

type TypeWriter func(s *typeSession, t reflect.Type, req *http.Request) error

type TypeServer struct {
	feed  <-chan interface{}
	write TypeWriter
}

type typeSession struct {
	depth int
	buf   *bufio.Writer
}

func (server TypeServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session := &typeSession{
		depth: 0,
		buf:   bufio.NewWriter(resp),
	}
	readType := reflect.TypeOf(<-server.feed)
	err := server.write(session, readType, req)
	if err != nil {
		panic(err)
	}
	session.buf.Flush()
}

func ServeTypeViewer(addr string, inchan <-chan interface{}) {
	server := TypeServer{feed: inchan, write: HTMLTypeWriter}
	err := http.ListenAndServe(addr, server)
	if err != nil {
		panic(err)
	}
}

const css = `
div[data-kind] {
	position: relative;
	border: 1px solid red;
	padding: 0.2em;
	background-color: #eee
}
div[data-kind=ptr] {
	background-color: #ccc
}
div[data-kind=array] {
	background-color: #c7c7f7
}
div[data-kind=slice] {
	background-color: #ccccff
}
div[data-kind=chan] {
	background-color: #fcc
}
div[data-kind=map] {
	background-color: #cfc
}
div[data-kind=func] {
	background-color: #ffc
}
div[data-kind=interface] {
	background-color: #cff
}
div[data-kind=struct] {
	background-color: #ddd
}
div[data-kind]::before {
	content: attr(data-kind);
	position: relative;
}`

func HTMLTypeWriter(session *typeSession, t reflect.Type, req *http.Request) error {
	const submit = `<form method="post"><button type="submit">Next</button></form>`
	if req.Method != "POST" {
		// serve form on GET requests so favicon.ico and co don't skip object under inspection
		session.buf.WriteString(`<!DOCTYPE html><html><body>` + submit + `</body></html>`)
		return nil
	}
	// write leading...
	session.buf.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html><head><title>Go: '%s'</title><style>
`+css+`
</style></head><body>`+submit, t))
	// ignore errors for the calls; we can't reasonably handle them unless we add a buffer
	_ = mirror.Walk(t, session.typeToHTML)
	// close all tags
	_ = session.typeToHTML(nil, 0, 0)
	_ = session.closeHtmlTagsToDepth(0)
	// write closing code...
	_, _ = session.buf.WriteString(`</body></html>`)
	return nil
}

func (session *typeSession) closeHtmlTagsToDepth(depth int) error {
	for d := session.depth - depth; d > 0; d-- {
		_, err := session.buf.WriteString("</div>")
		if err != nil {
			return err
		}
	}
	return nil
}

func (session *typeSession) typeToHTML(t *reflect.StructField, typeIndex, depth int) error {
	// for now, we are error-ignorant
	_ = session.closeHtmlTagsToDepth(depth - 1)
	if t == nil {
		return nil
	}
	tt := t.Type
	_, _ = session.buf.WriteString(fmt.Sprintf(
		`<div data-kind="%s"`, tt.Kind()))
	if len(t.Index) > 0 {
		_, _ = session.buf.WriteString(fmt.Sprintf(
			` data-field="%s" data-index="%v" data-offset="%d" data-tag="%s" `,
			t.Name, t.Index, t.Offset, t.Tag))
	}
	switch tt.Kind() {
	case reflect.Chan:
		var direction string
		switch tt.ChanDir() {
		case reflect.RecvDir:
			direction = "receive"
		case reflect.SendDir:
			direction = "send"
		case reflect.BothDir:
			direction = "both"
		}
		session.buf.WriteString(` data-direction="` + direction + `"`)
	case reflect.Map:
		_, _ = session.buf.WriteString(fmt.Sprintf(
			` data-keytype="%s"`, tt.Key()))
	case reflect.Array:
		_, _ = session.buf.WriteString(fmt.Sprintf(
			` data-length="%d"`, tt.Len()))
	case reflect.Func:
		_, _ = session.buf.WriteString(fmt.Sprintf(
			` data-args-in="%d" data-args-out="%d"`, tt.NumIn(), tt.NumOut()))
	}
	_, _ = session.buf.WriteString(fmt.Sprintf(
		` data-type="%v" data-size="%d" data-typeid="%d">`,
		tt, tt.Size(), typeIndex))
	session.depth = depth
	return nil
}
