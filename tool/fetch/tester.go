package fetch

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type testFetcher map[string]*TestResponse

func Test(m map[string]*TestResponse) Fetcher {
	f := make(testFetcher, len(m))
	for k, r := range m {
		if !strings.Contains(k, "\n") {
			k = "GET\n" + k + "\n\n\n"
		}
		h := sha256.Sum256([]byte(k))
		f[hex.EncodeToString(h[:])] = r
	}
	return f
}

func (testFetcher) Name() string { return "test" }

func (tf testFetcher) Fetch(request *Request) (*Response, error) {
	r, ok := tf[request.ID()]
	if !ok {
		return nil, fmt.Errorf("Not found %s %s", request.Method, request.URL.String())
	}

	return r.Response(), nil
}

type TestResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

// Create a test response.
// Headers is key1, value1, key2, value2 ...
// If headers length is odd, the last element is ignored.
func TR(status int, body []byte, headers ...string) *TestResponse {
	h := make(http.Header, len(headers)/2)
	for i := 0; i+1 < len(headers); i += 2 {
		h.Add(headers[i], headers[i+1])
	}
	return &TestResponse{
		Status: status,
		Header: h,
		Body:   body,
	}
}

// Like [TR] but the body is a string
func TRs(status int, body string, headers ...string) *TestResponse {
	return TR(status, []byte(body), headers...)
}

func (r *TestResponse) Response() *Response {
	return &Response{
		Status: r.Status,
		Header: r.Header.Clone(),
		Body:   io.NopCloser(bytes.NewReader(r.Body)),
	}
}
