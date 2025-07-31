package fetch

import (
	"bytes"
	"fmt"
	"io"
)

type testFetcher map[string]*Meta

// Create a Fetcher for testing.
// Use [Request.T] or [Request.Ts] functions for create [Meta] struct.
func Test(m ...*Meta) Fetcher {
	f := make(testFetcher, len(m))
	for _, m := range m {
		f[m.ID()] = m
	}
	return f
}

func (testFetcher) Name() string { return "test" }

func (tf testFetcher) Fetch(request *Request) (*Response, error) {
	m := tf[request.ID()]
	if m == nil {
		return nil, fmt.Errorf("Not found %s %s", request.Method, request.URL.String())
	}
	return &Response{
		Status: m.Status,
		Header: m.ResponseHeader.Clone(),
		Body:   io.NopCloser(bytes.NewReader(m.ResponseBody)),
	}, nil
}

// Create a [Meta] struct to use with [Test].
// Headers is key1, value1, key2, value2 ...
// If headers length is odd, the last element is ignored.
func (r *Request) T(status int, body []byte, headers ...string) *Meta {
	return &Meta{
		Method:        r.Method,
		URL:           r.URL,
		RequestHeader: r.Header,
		RequestBody:   r.Body,

		Status:         status,
		ResponseHeader: makeHeaders(headers),
		ResponseBody:   body,
	}
}

// Like [Request.T] with a string body.
func (r *Request) Ts(status int, body string, headers ...string) *Meta {
	return r.T(status, []byte(body), headers...)
}
