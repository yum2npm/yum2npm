package utils

import (
	"io"
	"log/slog"
	"net/http"
	"time"
)

func FetchUrl(url string) ([]byte, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	if res, err := client.Get(url); err != nil {
		return nil, err
	} else {
		defer func() {
			if err := res.Body.Close(); err != nil {
				slog.Error("Failed to close response body", "Error", err)
			}
		}()

		return io.ReadAll(res.Body)
	}
}
