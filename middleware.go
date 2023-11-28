package scihttp

import (
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
)

type MiddlewareChainNode interface {
	Next(m *Middleware) MiddlewareChainNode
	Handler(h httprouter.Handle) httprouter.Handle
}

type Middleware struct {
	pipelineRoot *MiddlewarePipeline
	preProcess   http.HandlerFunc
	postProcess  http.HandlerFunc
	next         *Middleware
	handler      httprouter.Handle
}

func (mdw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := w.(*httptest.ResponseRecorder)
	mdw.preProcess(rec, r)
	if mdw.next != nil {
		mdw.next.ServeHTTP(rec, r)
	} else if (mdw.handler != nil) && (rec.Code == 200) {
		params := httprouter.ParamsFromContext(r.Context())
		mdw.handler(rec, r, params)
	}
	mdw.postProcess(rec, r)
}

func (mdw *Middleware) Next(m *Middleware) MiddlewareChainNode {
	mdw.next = m
	m.pipelineRoot = mdw.pipelineRoot
	return mdw.next
}

func (mdw *Middleware) Handler(h httprouter.Handle) httprouter.Handle {
	mdw.handler = h
	return mdw.pipelineRoot.HandleHTTP
}

type MiddlewarePipeline struct {
	firstMiddleware *Middleware
	handler         httprouter.Handle
}

func (mp *MiddlewarePipeline) HandleHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	rec := httptest.NewRecorder()
	if mp.firstMiddleware != nil {
		mp.firstMiddleware.ServeHTTP(rec, r)
	} else if mp.handler != nil {
		params := httprouter.ParamsFromContext(r.Context())
		mp.handler(rec, r, params)
	}
	for k, v := range rec.Header() {
		w.Header()[k] = v
	}
	w.WriteHeader(rec.Result().StatusCode)
	w.Write(rec.Body.Bytes())
}

func (mp *MiddlewarePipeline) Next(m *Middleware) MiddlewareChainNode {
	mp.firstMiddleware = m
	m.pipelineRoot = mp
	return mp.firstMiddleware
}

func (mp *MiddlewarePipeline) Handler(h httprouter.Handle) httprouter.Handle {
	mp.handler = h
	return mp.HandleHTTP
}

func NewMiddleware(pre http.HandlerFunc, post http.HandlerFunc) *Middleware {
	return &Middleware{preProcess: pre, postProcess: post}
}

func NewMiddlewarePipeline() MiddlewareChainNode {
	mp := &MiddlewarePipeline{}
	return mp
}
