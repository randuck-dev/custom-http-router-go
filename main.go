package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"time"

	api "github.com/randuck-dev/rd-api/pkg"
)

func main() {
	var routerType string
	flag.StringVar(&routerType, "router", "custom", "router type")

	if routerType == "custom" {
		customRouter()
	}
	if routerType == "standard" {
		standardLibraryHttpHandler()
	}

	panic("Invalid router type")
}

// This spins up a simple custom router where handlers can be registered and the context can be extended
func customRouter() {

	ca := api.CustomerApi{}
	oa := api.OrderApi{}

	db := api.NewCustomerDatabase()
	defer db.Close()

	healthzApi := api.NewHealthzApi(db)

	executor := api.Executor{
		Name: "AGGRESSOR",
	}

	router := api.NewRouter(&executor)

	router.Handler("/customer", ca)
	router.Handler("/order", oa)
	router.Handler("/healthz", healthzApi)

	server := &http.Server{
		Addr: ":9092",
		BaseContext: func(l net.Listener) context.Context {
			ctx := context.Background()
			ctx = api.NewExecutorContext(ctx, &executor)
			ctx = api.WithDbContext(ctx, db)
			return ctx
		},

		Handler: &router,
	}

	// go dbStats(db)

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("[CUSTOM_SERVER] Something happened", "err", err)
	}
}

func dbStats(db *sql.DB) {
	// every 5 seconds, show the stats for the database

	for {
		stats := db.Stats()

		slog.Info("DB Stats", "stats", stats)
		time.Sleep(5 * time.Second)
	}

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
