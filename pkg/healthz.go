package api

import (
	"database/sql"
	"log/slog"
	"net/http"
)

type HealthzApi struct {
	db *sql.DB
}

func NewHealthzApi(db *sql.DB) HealthzApi {
	return HealthzApi{db: db}
}

func (ha HealthzApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := ha.db.Ping()
	if err != nil {
		slog.Error("Unable to ping database", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	okJsonResponse(w, []byte(`{"status": "ok"}`))
}
