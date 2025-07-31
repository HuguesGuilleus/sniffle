package fetch_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"sniffle/tool/fetch"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveHTTP(t *testing.T) {
	u, err := url.Parse("https://example.com/file?b=2&a=1")
	assert.NoError(t, err)
	request := fetch.Request{
		Method: http.MethodPost,
		URL:    u,
		Header: http.Header{
			"accept-encoding": []string{"application/json"},
		},
		Body: []byte("yo!"),
	}
	response := fetch.Response{
		Status: 200,
		Header: http.Header{
			"content-type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte("Hello World!"))),
	}
	now := time.Date(2024, time.November, 6, 15, 46, 05, 0, time.UTC)

	// Save
	buff := bytes.Buffer{}
	assert.NoError(t, fetch.SaveHTTP(&request, &response, now, &buff))

	// Read
	savedResponse, err := fetch.ReadResponse(io.NopCloser(&buff))
	assert.NoError(t, err)
	assert.Equal(t, &fetch.Response{
		Status: 200,
		Header: http.Header{
			"content-type": []string{"application/json"},
		},
		Body: io.NopCloser(&buff),
	}, savedResponse)

	data, err := io.ReadAll(savedResponse.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Hello World!", string(data))

	// ReadMeta
	buff.Reset()
	assert.NoError(t, fetch.SaveHTTP(&request, &response, now, &buff))
	m, err := fetch.ReadMeta(io.NopCloser(&buff))
	assert.NoError(t, err)
	assert.Equal(t, &fetch.Meta{
		Time:   now,
		Method: http.MethodPost,
		RawURL: u.String(),
		URL:    u,
		RequestHeader: http.Header{
			"accept-encoding": []string{"application/json"},
		},
		RequestBody: []byte("yo!"),
		Status:      200,
		ResponseHeader: http.Header{
			"content-type": []string{"application/json"},
		},
	}, m)

	// multiple call to check cache of id and path.
	assert.Equal(t, "e3f2c2ce79cb587525022983d1229bed07fbb499ff9efe02de768e2f173379bc", m.ID())
	assert.Equal(t, "e3f2c2ce79cb587525022983d1229bed07fbb499ff9efe02de768e2f173379bc", m.ID())
	assert.Equal(t, "https/example.com/e3f2c2ce79cb587525022983d1229bed07fbb499ff9efe02de768e2f173379bc.http", m.Path())
	assert.Equal(t, "https/example.com/e3f2c2ce79cb587525022983d1229bed07fbb499ff9efe02de768e2f173379bc.http", m.Path())
}
