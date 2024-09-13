package fetch

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNet(t *testing.T) {
	assert.Equal(t, "net", Net(nil, "", 0, 0).Name())

	u, err := url.Parse("https://example.com/file?b=2&a=1")
	assert.NoError(t, err)

	body, id, err := Net(fakeRoundTrip{}, "cache", 0, 0).Fetch(
		context.Background(),
		"PATCH",
		u,
		http.Header{"H": []string{"42"}},
		[]byte("body"))
	assert.NoError(t, err)
	assert.Equal(t, "7be54c1355", id)
	defer body.Close()

	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.EqualValues(t, "body", data)
	assert.IsType(t, &os.File{}, body)

	assert.NoError(t, os.RemoveAll("cache"))
}

type fakeRoundTrip struct{}

func (fakeRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Method != "PATCH" {
		panic("not PATCH method")
	}
	if r.Header.Get("H") != "42" {
		panic("wrong headear")
	}
	if body, err := io.ReadAll(r.Body); err != nil {
		panic("fail to read body: " + err.Error())
	} else if string(body) != "body" {
		panic("wrong body")
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"yolo": []string{"v1", "v2"}},
		Body:       io.NopCloser(bytes.NewReader([]byte("body"))),
	}, nil
}
