package fetch

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	u, err := url.Parse("https://example.com/file?b=2&a=1")
	assert.NoError(t, err)

	r := &Request{
		URL: u,
		Header: http.Header{
			"accept-encoding": []string{"application/json"},
		},
		Body: []byte("Hello World!"),
	}

	id := "150ee78242120c0e38fc747a175c56068c2f07f8b0c57345a7ee6cdd5a172d05"
	assert.Equal(t, "https/example.com/150ee78242120c0e38fc747a175c56068c2f07f8b0c57345a7ee6cdd5a172d05.http", r.Path())
	assert.Equal(t, id, r.id)
	assert.Equal(t, id, r.ID())
	assert.Equal(t, "GET", r.Method)

	r = &Request{URL: u}
	r.ID()
	assert.Equal(t, "GET", r.Method)
	assert.NotNil(t, r.Header)

	assert.Equal(t, `err://http://example.com/%25%25`, URL("http://example.com/%%").URL.String())
}
