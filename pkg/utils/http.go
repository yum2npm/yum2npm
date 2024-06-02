package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Length", strconv.Itoa(len(error)))
	w.WriteHeader(code)
	fmt.Fprint(w, error)
}

func NotFound(w http.ResponseWriter) {
	Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func JsonResponse(w http.ResponseWriter, data any) {
	j, err := json.Marshal(data)
	if err != nil {
		Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Error("error marshalling json", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(j)))
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
