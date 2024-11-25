package fetch

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"maps"
	"net/http"
	"net/url"
	"path/filepath"
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

// Create a *Request from a URL.
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

// Easy way to create a request.
// If error, create a request with url: scheme=err, path=rawURL.
//
// Headers is key1, value1, key2, value2 ...
// If headers length is odd, the last element is ignored.
//
// Use .B() or Bs() to add a body.
func R(method, rawURL string, body []byte, headers ...string) *Request {
	r := URL(rawURL)
	r.Method = method
	r.Body = body
	r.Header = make(http.Header, len(headers)/2)
	for i := 0; i+1 < len(headers); i += 2 {
		r.Header.Add(headers[i], headers[i+1])
	}
	return r
}

// Like [R] with body as string.
func Rs(method, rawURL, body string, headers ...string) *Request {
	return R(method, rawURL, []byte(body), headers...)
}

// Cononize the request and return is ID.
// The ID is a hash of all fields.
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
	for k, v := range r.Header {
		k = http.CanonicalHeaderKey(k)
		for _, v := range v {
			h.Write([]byte(k))
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

func getPath(base string, request *Request) string {
	return filepath.Join(
		base,
		request.URL.Scheme,
		request.URL.Host,
		request.ID()+".http",
	)
}

func getDir(base string, request *Request) string {
	return filepath.Join(
		base,
		request.URL.Scheme,
		request.URL.Host,
	)
}
