package router

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type routerHandler struct {
	pattern string
	handler http.Handler
}

type middleware func(http.Handler) http.Handler

type Router struct {
	handlers    map[string]routerHandler
	middlewares []middleware
}

func NewRouter() Router {
	return Router{
		handlers:    make(map[string]routerHandler),
		middlewares: []middleware{},
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

func (router *Router) Middleware(handler middleware) {
	router.middlewares = append(router.middlewares, handler)
}

func middlewareChain(middlewares []middleware, handler http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := r.URL.Path

	var handler http.Handler

	// Might be lucky and have an exact match
	if rh, ok := router.handlers[requestPath]; ok {
		handler = rh.handler
	} else {
		for prefix, rh := range router.handlers {
			if strings.HasPrefix(requestPath, prefix) {
				handler = rh.handler
				break
			}
		}
	}

	if handler == nil {
		slog.Error("No handler found", "path", requestPath)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	handler = middlewareChain(router.middlewares, handler)

	handler.ServeHTTP(w, r)

}

func RequestDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start)

		slog.Info("Request", "path", r.URL.Path, "method", r.Method, "elapsed", elapsed)
	})
}
