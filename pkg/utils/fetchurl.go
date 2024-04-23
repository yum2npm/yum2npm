package utils

import (
	"io"
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
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}
}
