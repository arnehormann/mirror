package httpexport

import (
	"bufio"
	"fmt"
	"github.com/arnehormann/mirror"
	"io"
	"net/http"
	"reflect"
)

func NewTypeServer(addr string) chan<- interface{} {
	typechan := make(chan interface{})
	go func(addr string, inchan <-chan interface{}) {
		server := typeServer{feed: inchan, write: htmlTypeWriter}
		err := http.ListenAndServe(addr, server)
		if err != nil {
			panic(err)
		}
	}(addr, typechan)
	return typechan
}

type typeWriter func(s *typeSession, t *reflect.Type) error

type typeServer struct {
	feed  <-chan interface{}
	write typeWriter
}

type typeSession struct {
	depth int
	buf   *bufio.Writer
}

func (server typeServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	session := &typeSession{
		depth: 0,
		buf:   bufio.NewWriter(resp),
	}
	var t *reflect.Type
	if req.Method == "POST" {
		readType := reflect.TypeOf(<-server.feed)
		t = &readType
	}
	must(server.write(session, t))
	session.buf.Flush()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func mustIO(processed int, err error) {
	switch {
	case err == nil, err == io.EOF:
		return
	}
	panic(err)
}

// code for html type export

func htmlTypeWriter(session *typeSession, t *reflect.Type) error {
	const submit = `<form method="post"><button type="submit">Next</button></form>`
	if t == nil {
		// serve form on GET requests so favicon.ico and co don't skip object under inspection
		mustIO(session.buf.WriteString(`<!DOCTYPE html><html><body>` + submit + `</body></html>`))
		return nil
	}
	// write leading...
	mustIO(session.buf.WriteString(fmt.Sprintf(`
<!DOCTYPE html>
<html><head><title>Go: '%s'</title><style>
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
	content: attr(data-kind) ': ' attr(data-field) ' ' attr(data-type);
	position: relative;
	margin-left: 1em;
}
</style></head><body>%s`, t, submit)))
	// ignore errors for the calls; we can't reasonably handle them unless we add a buffer
	must(mirror.Walk(*t, session.typeToHTML))
	// close all tags
	must(session.typeToHTML(nil, 0, 0))
	must(session.closeHtmlTagsToDepth(0))
	// write closing code...
	mustIO(session.buf.WriteString(`</body></html>`))
	return nil
}

func (session *typeSession) closeHtmlTagsToDepth(depth int) error {
	for d := session.depth - depth; d >= 0; d-- {
		_, err := session.buf.WriteString("</div>")
		if err != nil {
			return err
		}
	}
	return nil
}

func (session *typeSession) typeToHTML(t *reflect.StructField, typeIndex, depth int) error {
	// for now, we are error-ignorant
	must(session.closeHtmlTagsToDepth(depth))
	if t == nil {
		return nil
	}
	tt := t.Type
	mustIO(session.buf.WriteString(fmt.Sprintf(
		`<div data-kind="%s" data-type="%s" data-size="%d" data-typeid="%d"`,
		tt.Kind(), tt, tt.Size(), typeIndex)))
	if len(t.Index) > 0 {
		mustIO(session.buf.WriteString(fmt.Sprintf(
			` data-field="%s" data-index="%v" data-offset="%d" data-tag="%s" `,
			t.Name, t.Index, t.Offset, t.Tag)))
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
		mustIO(session.buf.WriteString(fmt.Sprintf(
			` data-keytype="%s"`, tt.Key())))
	case reflect.Array:
		mustIO(session.buf.WriteString(fmt.Sprintf(
			` data-length="%d"`, tt.Len())))
	case reflect.Func:
		mustIO(session.buf.WriteString(fmt.Sprintf(
			` data-args-in="%d" data-args-out="%d"`, tt.NumIn(), tt.NumOut())))
	}
	mustIO(session.buf.WriteString(`>`))
	session.depth = depth
	return nil
}
