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

	body, id, err := Net(fakeRoudnTrip{}, "cache", 0, 0).Fetch(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, "62caa69659", id)
	defer body.Close()

	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.EqualValues(t, "body", data)
	assert.IsType(t, &os.File{}, body)

	assert.NoError(t, os.RemoveAll("cache"))
}

type fakeRoudnTrip struct{}

func (fakeRoudnTrip) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"yolo": []string{"v1", "v2"}},
		Body:       io.NopCloser(bytes.NewReader([]byte("body"))),
	}, nil
}
