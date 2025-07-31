package fetch

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Information about a request and response, saved in cache file or used for the test.
type Meta struct {
	Time time.Time `json:"time"`

	Method        string      `json:"requestMethod"`
	RawURL        string      `json:"requestURL"`
	URL           *url.URL    `json:"-"`
	RequestHeader http.Header `json:"requestHeader"`
	RequestBody   []byte      `json:"requestBody"`

	Status         int         `json:"status"`
	ResponseHeader http.Header `json:"responseHeader"`

	// Response body used only for test
	ResponseBody []byte `json:"-"`

	id   string `json:"-"`
	path string `json:"-"`
}

// Calculate the ID with [Request.ID].
// Cache the id, so you must not edit m after this method call.
func (m *Meta) ID() string {
	if m.id != "" {
		return m.id
	}

	m.id = (&Request{
		Method: m.Method,
		URL:    m.URL,
		Header: m.RequestHeader,
		Body:   m.RequestBody,
	}).ID()

	return m.id
}

// Return the path where save the request.
// Use [Meta.ID] method, so you must not edit m after this method call.
func (m *Meta) Path() string {
	if m.path != "" {
		return m.path
	}
	m.path = getPath("", &Request{
		URL: m.URL,
		id:  m.ID(),
	})
	return m.path
}

// Write meta information then response body into w writer.
// See README.md for format details.
func SaveHTTP(request *Request, response *Response, now time.Time, w io.Writer) error {
	defer response.Body.Close()

	buff := bytes.Buffer{}
	buff.WriteString("HTTP\x00\x00\x00\x00")
	json.NewEncoder(&buff).Encode(&Meta{
		Time: now.UTC(),

		Method:        request.Method,
		RawURL:        request.URL.String(),
		RequestHeader: request.Header,
		RequestBody:   request.Body,

		Status:         response.Status,
		ResponseHeader: response.Header,
	})
	buff.WriteString("\n\n")
	binary.BigEndian.PutUint32(buff.Bytes()[4:8], uint32(buff.Len()-8))

	if _, err := w.Write(buff.Bytes()); err != nil {
		return err
	}
	if _, err := io.Copy(w, response.Body); err != nil {
		return err
	}

	return nil
}

func ReadResponse(r io.ReadCloser) (*Response, error) {
	meta, err := ReadOnlyMeta(r)
	if err != nil {
		r.Close()
		return nil, err
	}

	return &Response{
		Status: meta.Status,
		Header: meta.ResponseHeader,
		Body:   r,
	}, nil
}

// Read meta data and close the reader.
func ReadMeta(r io.ReadCloser) (*Meta, error) {
	defer r.Close()
	return ReadOnlyMeta(r)
}

// Read only header of the file, then you can read the body.
func ReadOnlyMeta(r io.Reader) (*Meta, error) {
	head := [8]byte{}
	if _, err := io.ReadFull(r, head[:]); err != nil {
		return nil, err
	} else if !bytes.Equal(head[0:4], []byte{'H', 'T', 'T', 'P'}) {
		return nil, fmt.Errorf("wrong head: %q", head[0:4])
	}

	j := make([]byte, binary.BigEndian.Uint32(head[4:8]))
	meta := &Meta{}
	if _, err := io.ReadFull(r, j); err != nil {
		return nil, err
	} else if err := json.Unmarshal(j, meta); err != nil {
		return nil, err
	}

	u, err := url.Parse(meta.RawURL)
	if err != nil {
		return nil, err
	}
	meta.URL = u

	return meta, nil
}
