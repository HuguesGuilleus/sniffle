package fetch_test

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"sniffle/tool/fetch"
	"sniffle/tool/writefs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClearCache(t *testing.T) {
	defer os.RemoveAll("_cache")

	now := time.Now().UTC()
	request := fetch.Rs("", "https://cache.net/dir/file.txt?a=1", "body", "k1", "v1", "k2", "v2")
	assert.NoError(t, os.MkdirAll("_cache/https/cache.net/", 0o775))
	path := "_cache/" + request.Path()
	f, err := os.Create(path)
	assert.NoError(t, err)
	fetch.SaveHTTP(
		request,
		&fetch.Response{
			Status: 200,
			Header: http.Header{
				"H1": []string{"v1", "v2"},
			},
			Body: io.NopCloser(bytes.NewReader([]byte("..."))),
		},
		now,
		f,
	)
	assert.NoError(t, f.Close())

	call := 0
	assert.NoError(t, fetch.ClearCache(writefs.Os("_cache"), func(m *fetch.Meta) time.Duration {
		call++
		assert.Equal(t, &fetch.Meta{
			Time:   now,
			Method: "GET",
			RawURL: "https://cache.net/dir/file.txt?a=1",
			URL:    request.URL,
			RequestHeader: http.Header{
				"K1": []string{"v1"},
				"K2": []string{"v2"},
			},
			RequestBody: []byte("body"),
			Status:      200,
			ResponseHeader: http.Header{
				"H1": []string{"v1", "v2"},
			},
		}, m)
		return 0
	}))
	assert.Equal(t, 1, call)
	assert.NoFileExists(t, path)
}
