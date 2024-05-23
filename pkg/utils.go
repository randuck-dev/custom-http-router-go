package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func okJsonResponse(w http.ResponseWriter, payload []byte) {
	w.Header().Add("Content-Type", "application/json")

	slog.Info("Sending payload", "content-length", len(payload))
	w.Write(payload)
}

func errorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)

	errResponse := Error{
		errorMsg:  "<<< TODO >>>",
		errorCode: http.StatusBadRequest,
	}

	serialized, err := json.Marshal(errResponse)

	if err != nil {
		slog.Error("Unable to create error response", "error", err)
		return
	}
	okJsonResponse(w, serialized)
}

func internalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

type Error struct {
	errorMsg  string `json:"error"`
	errorCode uint32 `json:"errorCode"`
}
