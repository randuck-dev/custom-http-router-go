package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
)

type Router struct {
	executor *Executor
	handlers map[string]http.Handler
}

func NewRouter(executor *Executor) Router {
	return Router{
		executor: executor,
		handlers: make(map[string]http.Handler),
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

	router.handlers[path] = handler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Debug("Received request", "path", r.URL.Path)

	for path, handler := range router.handlers {
		if strings.HasPrefix(r.URL.Path, path) {
			handler.ServeHTTP(w, r)
		}
	}
}
