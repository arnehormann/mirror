package httpexport

import (
	"bufio"
	"fmt"
	"github.com/arnehormann/mirror"
	"net/http"
	"reflect"
	"strings"
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
	err   error
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
	err := server.write(session, t)
	if err != nil {
		panic(err)
	}
	session.buf.Flush()
}

func (session *typeSession) Concat(text string) {
	if session.err != nil {
		return
	}
	_, err := session.buf.WriteString(text)
	session.err = err
}

func (session *typeSession) Concatf(format string, args ...interface{}) {
	if session.err != nil {
		return
	}
	_, err := session.buf.WriteString(fmt.Sprintf(format, args...))
	session.err = err
}

// code for html type export

func htmlTypeWriter(session *typeSession, t *reflect.Type) error {
	const submit = `<form method="post"><button type="submit">Next</button></form>`
	if t == nil {
		// serve form on GET requests so favicon.ico and co don't skip object under inspection
		session.Concat(`<!DOCTYPE html><html><body>` + submit + `</body></html>`)
		return session.err
	}
	// write leading...
	session.Concatf(`
<!DOCTYPE html>
<html><head><title>Go: '%s'</title><style>
div[data-kind] {
	position: relative;
	border: 1px solid red;
	padding: 0.2em;
	background-color: #eee
}
div[data-kind]::before {
	content: attr(data-kind) ': ' attr(data-field) ' ' attr(data-type);
	position: relative;
	margin-left: 1em;
}
div[data-kind=ptr]			{ background-color: #cccccc }
div[data-kind=array]		{ background-color: #c7c7f7 }
div[data-kind=slice]		{ background-color: #ccccff }
div[data-kind=chan]			{ background-color: #ffcccc }
div[data-kind=map]			{ background-color: #ccffcc }
div[data-kind=func]			{ background-color: #ffffcc }
div[data-kind=interface]	{ background-color: #ccffff }
div[data-kind=struct]		{ background-color: #dddddd }
</style></head><body>%s`, *t, submit)
	typeToHtml := func(t *reflect.StructField, typeIndex, depth int) error {
		// for now, we are error-ignorant
		// close open tags
		if session.depth > depth {
			session.Concat(strings.Repeat("</div>", session.depth-depth))
		}
		// if no type is given, return
		if t == nil {
			return nil
		}
		tt := t.Type
		session.Concatf(
			`<div data-kind="%s" data-type="%s" data-size="%d" data-typeid="%d"`,
			tt.Kind(), tt, tt.Size(), typeIndex)
		if len(t.Index) > 0 {
			session.Concatf(
				` data-field="%s" data-index="%v" data-offset="%d" data-tag="%s" `,
				t.Name, t.Index, t.Offset, t.Tag)
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
			session.Concat(` data-direction="` + direction + `"`)
		case reflect.Map:
			session.Concatf(` data-keytype="%s"`, tt.Key())
		case reflect.Array:
			session.Concatf(` data-length="%d"`, tt.Len())
		case reflect.Func:
			session.Concatf(` data-args-in="%d" data-args-out="%d"`, tt.NumIn(), tt.NumOut())
		}
		session.Concat(`>`)
		session.depth = depth
		return session.err
	}
	// ignore errors for the calls; we can't reasonably handle them unless we add a buffer
	session.err = mirror.Walk(*t, typeToHtml)
	if session.err != nil {
		return session.err
	}
	// close all tags
	session.err = typeToHtml(nil, 0, 0)
	// write closing code...
	session.Concat(`</body></html>`)
	return session.err
}
