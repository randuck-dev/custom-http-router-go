package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"

	api "github.com/randuck-dev/rd-api/pkg"
)

func main() {
	customRouter()
	standardLibraryHttpHandler()
}

// This spins up a simple custom router where handlers can be registered and the context can be extended
func customRouter() {
	ca := api.CustomerApi{}
	oa := api.OrderApi{}

	executor := api.Executor{
		Name: "AGGRESSOR",
	}

	router := api.NewRouter(&executor)

	router.Handler("/customer", ca)
	router.Handler("/order", oa)

	server := &http.Server{
		Addr: ":9092",
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.Background()
			ctx = api.NewExecutorContext(ctx, &executor)
			return ctx
		},

		Handler: &router,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("[CUSTOM_SERVER] Something happened", "err", err)
		}
	}()
}

func standardLibraryHttpHandler() {
	ca := api.CustomerApi{}
	oa := api.OrderApi{}
	http.HandleFunc("/customer", ca.ServeHTTP)
	http.HandleFunc("/order", oa.ServeHTTP)

	err := http.ListenAndServe(":9091", nil)

	if err != nil {
		slog.Error("[http] Something happened", "err", err)
	}

}
