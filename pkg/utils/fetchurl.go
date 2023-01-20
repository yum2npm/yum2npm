package utils

import (
	"io"
	"net/http"
)

func FetchUrl(url string) ([]byte, error) {
	if res, err := http.Get(url); err != nil {
		return make([]byte, 0), err
	} else {
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}
}
