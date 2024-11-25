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

type Meta struct {
	Time time.Time `json:"time"`

	Method        string      `json:"requestMethod"`
	RawURL        string      `json:"requestURL"`
	URL           *url.URL    `json:"-"`
	RequestHeader http.Header `json:"requestHeader"`
	RequestBody   []byte      `json:"requestBody"`

	Status         int         `json:"status"`
	ResponseHeader http.Header `json:"responseHeader"`
}

func SaveHTTP(request *Request, response *Response, now time.Time, w io.Writer) error {
	j, _ := json.Marshal(&Meta{
		Time: now.UTC(),

		Method:        request.Method,
		RawURL:        request.URL.String(),
		RequestHeader: request.Header,
		RequestBody:   request.Body,

		Status:         response.Status,
		ResponseHeader: response.Header,
	})

	head := [8]byte{'H', 'T', 'T', 'P'}
	binary.BigEndian.PutUint32(head[4:], uint32(len(j)+2))

	if _, err := w.Write(head[:]); err != nil {
		return err
	}
	if _, err := w.Write(j); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'\n', '\n'}); err != nil {
		return err
	}
	if _, err := io.Copy(w, response.Body); err != nil {
		return err
	}

	return response.Body.Close()
}

func ReadResponse(r io.ReadCloser) (*Response, error) {
	needClose := true
	defer func() {
		if needClose {
			r.Close()
		}
	}()

	head := [8]byte{}
	if _, err := r.Read(head[:]); err != nil {
		return nil, err
	} else if !bytes.Equal(head[0:4], []byte{'H', 'T', 'T', 'P'}) {
		return nil, fmt.Errorf("wrong head: %q", head[0:4])
	}

	j := make([]byte, binary.BigEndian.Uint32(head[4:8]))
	meta := Meta{}
	if _, err := r.Read(j); err != nil {
		return nil, err
	} else if err := json.Unmarshal(j, &meta); err != nil {
		return nil, err
	}

	needClose = false
	return &Response{
		Status: meta.Status,
		Header: meta.ResponseHeader,
		Body:   r,
	}, nil
}

// Read meta data and close the reader.
func ReadMeta(r io.ReadCloser) (*Meta, error) {
	defer r.Close()

	head := [8]byte{}
	if _, err := r.Read(head[:]); err != nil {
		return nil, err
	} else if !bytes.Equal(head[0:4], []byte{'H', 'T', 'T', 'P'}) {
		return nil, fmt.Errorf("wrong head: %q", head[0:4])
	}

	j := make([]byte, binary.BigEndian.Uint32(head[4:8]))
	meta := &Meta{}
	if _, err := r.Read(j); err != nil {
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
