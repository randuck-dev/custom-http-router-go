package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type routerHandler struct {
	pattern string
	handler http.Handler
}

type Router struct {
	executor *Executor
	handlers map[string]routerHandler
}

func NewRouter(executor *Executor) Router {
	return Router{
		executor: executor,
		handlers: make(map[string]routerHandler),
	}
}

func (router *Router) ExtendContext(f func(context.Context)) {
	ctx := context.Background()
	f(ctx)
}

func (router *Router) Handler(path string, handler http.Handler) {
	if _, ok := router.handlers[path]; ok {
		slog.Error("Handler already registered!", "path", path)
		panic("Handler already registered")
	}

	router.handlers[path] = routerHandler{
		pattern: path,
		handler: handler,
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path

	// Might be lucky and have an exact match
	if rh, ok := router.handlers[requestPath]; ok {
		rh.handler.ServeHTTP(w, r)
		return
	}

	for prefix, rh := range router.handlers {
		if strings.HasPrefix(requestPath, prefix) {
			rh.handler.ServeHTTP(w, r)
		}
	}
}
