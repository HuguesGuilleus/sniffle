package fetch_test

import (
	"bytes"
	"io"
	"net/http"
	"sniffle/tool/fetch"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester(t *testing.T) {
	f := fetch.Test(
		fetch.Fmt("https://example.org/dir/%s", "file").Ts(200, "Hello World!", "Content-Type", "text/plain; charset=utf-8"),
		fetch.Rs("PATCH", "https://example.org/api/title", `{"title":"Yolo"}`, "Content-Type", "application/json").T(205, nil, "X-Header"),
	)
	assert.Equal(t, "test", f.Name())

	response, err := f.Fetch(fetch.R("", "https://example.org/dir/file", nil))
	assert.NoError(t, err)
	assert.Equal(t, &fetch.Response{
		Status: 200,
		Header: http.Header{
			"Content-Type": []string{"text/plain; charset=utf-8"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte("Hello World!"))),
	}, response)

	response, err = f.Fetch(fetch.Rs("PATCH", "https://example.org/api/title",
		`{"title":"Yolo"}`,
		"Content-Type", "application/json",
	))
	assert.NoError(t, err)
	assert.Equal(t, &fetch.Response{
		Status: 205,
		Header: http.Header{},
		Body:   io.NopCloser(bytes.NewReader(nil)),
	}, response)

	response, err = f.Fetch(fetch.URL("https://example.org/404"))
	assert.Error(t, err)
	assert.Nil(t, response)
}
