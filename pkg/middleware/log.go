package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received := time.Now()
		newWriter := &ResponseWriter{w, http.StatusOK}
		next.ServeHTTP(newWriter, r)
		duration := time.Since(received)
		slog.Info(
			"Handled request",
			"Remote-Address", r.RemoteAddr,
			"Timestamp", received.Format(time.RFC3339),
			"Method", r.Method,
			"URI", r.RequestURI,
			"Duration", duration,
			"Content-Length", newWriter.Header().Get("Content-Length"),
			"Status", newWriter.statusCode,
			"User-Agent", r.UserAgent(),
		)
	})
}
