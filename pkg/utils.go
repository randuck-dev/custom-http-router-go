package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func okJsonResponse(w http.ResponseWriter, payload []byte) {
	w.Header().Add("Content-Type", "application/json")

	slog.Debug("Sending payload", "content-length", len(payload))
	if _, err := w.Write(payload); err != nil {
		slog.Error("Unable to write response", "error", err)
	}
}

func errorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)

	errResponse := Error{
		ErrorMsg:  "<<< TODO >>>",
		ErrorCode: http.StatusBadRequest,
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
	ErrorMsg  string `json:"error"`
	ErrorCode uint32 `json:"errorCode"`
}
