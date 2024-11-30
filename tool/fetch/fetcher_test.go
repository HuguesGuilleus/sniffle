package fetch_test

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"os"
	"sniffle/tool/fetch"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNetAndCache(t *testing.T) {
	defer os.RemoveAll("_cache")

	u, err := url.Parse("https://example.net/")
	assert.NoError(t, err)
	request := &fetch.Request{
		Method: http.MethodPatch,
		URL:    u,
		Header: http.Header{"H": []string{"42"}},
		Body:   []byte("body"),
	}

	testResponse := func(fetcher fetch.Fetcher) {
		response, err := fetcher.Fetch(request)
		assert.NoError(t, err)
		if err != nil {
			t.FailNow()
		}

		body := response.Body
		defer body.Close()
		response.Body = nil

		assert.Equal(t, &fetch.Response{
			Status: 200,
			Header: http.Header{"yolo": []string{"v1", "v2"}},
		}, response)

		data, err := io.ReadAll(body)
		assert.NoError(t, err)
		assert.Equal(t, "RESPONSE", string(data))
	}

	// Network
	assert.Equal(t, "net", fetch.Net(nil, "", nil).Name())
	testResponse(fetch.Net(fakeRoundTrip{}, "_cache", make(map[string]time.Duration)))

	// Cache
	assert.Equal(t, "cache", fetch.Cache("").Name())
	testResponse(fetch.Cache("_cache"))
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
		Body:       io.NopCloser(bytes.NewReader([]byte("RESPONSE"))),
	}, nil
}
