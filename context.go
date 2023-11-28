package scihttp

import (
	"context"
	"net/http"
)

type Context struct {
	context.Context
	request *http.Request
	writer  http.ResponseWriter
}

func (ctx *Context) ResponseWriter() http.ResponseWriter {
	return ctx.writer
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}
