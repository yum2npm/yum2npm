package utils

import (
	"bytes"
	"compress/gzip"
	"context"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/klauspost/compress/zstd"
	"github.com/ulikunitz/xz"
	"io"
	"net/http"
)

func FetchUrl(ctx context.Context, url string) (r io.Reader, err error) {
	client := http.DefaultClient

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	withContext := req.WithContext(ctx)

	var res *http.Response
	res, err = client.Do(withContext)
	if err != nil {
		return
	}

	var b []byte
	b, err = io.ReadAll(res.Body)

	var kind types.Type
	kind, err = filetype.Match(b)
	if err != nil {
		return
	}

	contentType := kind.MIME.Value

	if len(contentType) == 0 {
		contentType = res.Header.Get("Content-Type")
	}

	switch kind.MIME.Value {
	case "application/gzip":
		r, err = gzip.NewReader(bytes.NewReader(b))
	case "application/x-xz":
		r, err = xz.NewReader(bytes.NewReader(b))
	case "application/zstd":
		r, err = zstd.NewReader(bytes.NewReader(b))
	default:
		r = bytes.NewReader(b)
	}

	return
}
