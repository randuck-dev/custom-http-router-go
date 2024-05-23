package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type OrderApi struct{}

func (oa OrderApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "", "GET":
		oa.Get(w, r)
	default:
		slog.Info("Unrecognized method", "method", r.Method)
	}
}

func (oa OrderApi) Get(w http.ResponseWriter, r *http.Request) {

	payload := Order{
		Ticker:     "APPL@XNAS",
		LimitPrice: 20,
		Strategy:   AQUA,
	}

	serialized, err := json.Marshal(payload)
	if err != nil {
		errorResponse(w)
		return
	}

	okJsonResponse(w, serialized)
}

type Strategy int

const (
	AQUA    Strategy = iota
	DMA              = iota
	SORTDMA          = iota
)

type Order struct {
	Ticker     string
	LimitPrice float64
	Strategy   Strategy
}
