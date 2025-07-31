// Make HTTP requests.
package fetch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"path"
	"slices"
)

type Fetcher interface {
	// Make a HTTP GET request and return response.
	// If response code is out of [200,299] return an error.
	// Return fs.ErrNotExist if the fetcher is a cache, and not found the response.
	Fetch(*Request) (*Response, error)
	// Static name of the fetcher: net, cache...
	// Use this for the log
	Name() string
}

type Response struct {
	Status int // e.g. 200
	Header http.Header
	Body   io.ReadCloser
}

type Request struct {
	Method string
	URL    *url.URL
	Header http.Header

	Body []byte

	// evaluted ID
	id string
}

// Create a GET *Request from a URL.
// If error, create a request with url: scheme=err, path=rawURL.
func URL(rawURL string) *Request {
	u, _ := url.Parse(rawURL)
	if u == nil {
		u = new(url.URL)
		u.Scheme = "err"
		u.Path = rawURL
	}
	return &Request{URL: u}
}

// Like [URL] but call [fmt.Sprintf] to create the url.
func Fmt(format string, a ...any) *Request {
	return URL(fmt.Sprintf(format, a...))
}

// Easy way to create a request.
// If error, create a request with url: scheme=err, path=rawURL.
//
// Headers is key1, value1, key2, value2 ...
// If headers length is odd, the last element is ignored.
func R(method, rawURL string, body []byte, headers ...string) *Request {
	r := URL(rawURL)
	r.Method = method
	r.Body = body
	r.Header = makeHeaders(headers)
	return r
}

func makeHeaders(headers []string) http.Header {
	h := make(http.Header, len(headers)/2)
	for i := 0; i+1 < len(headers); i += 2 {
		h.Add(headers[i], headers[i+1])
	}
	return h
}

// Like [R] with body as string.
func Rs(method, rawURL, body string, headers ...string) *Request {
	return R(method, rawURL, []byte(body), headers...)
}

// Cononize the request and return the identifier.
// The ID is a hash of all fields encoded in hexadecimal.
//
// The ID is saved in the request, so after call to this method, you must not edit the request.
func (r *Request) ID() string {
	if r.id != "" {
		return r.id
	}

	if r.Method == "" {
		r.Method = http.MethodGet
	}
	r.URL.RawQuery = r.URL.Query().Encode()
	if r.Header == nil {
		r.Header = make(http.Header)
	}

	h := sha256.New()
	h.Write([]byte(r.Method))
	h.Write([]byte{'\n'})

	h.Write([]byte(r.URL.String()))
	h.Write([]byte{'\n'})

	keys := slices.Collect(maps.Keys(r.Header))
	slices.Sort(keys)
	for _, rawKey := range keys {
		canonicalKey := http.CanonicalHeaderKey(rawKey)
		for _, v := range r.Header[rawKey] {
			h.Write([]byte(canonicalKey))
			h.Write([]byte{':', ' '})
			h.Write([]byte(v))
			h.Write([]byte{'\n'})
		}
	}
	h.Write([]byte{'\n'})

	h.Write(r.Body)
	h.Write([]byte{'\n'})

	hash := [sha256.Size]byte{}
	h.Sum(hash[:0])

	r.id = hex.EncodeToString(hash[:])
	return r.id
}

func (r *Request) Path() string {
	return path.Join(
		r.URL.Scheme,
		r.URL.Host,
		r.ID()+".http",
	)
}
