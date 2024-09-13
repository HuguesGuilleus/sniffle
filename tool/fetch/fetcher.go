package fetch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
)

type Fetcher interface {
	// Make a HTTP GET request and return body with cache ID.
	// If response code is out of [200,299] return an error.
	Fetch(ctx context.Context, method string, u *url.URL, headers http.Header, body []byte) (io.ReadCloser, string, error)
	// Static name of the fetcher: net, cache...
	// Use this for the log
	Name() string
}

func GeneratePath(cacheBase, method string, u *url.URL, headers http.Header, body []byte) (logID, filePath string) {
	hasher := sha256.New()
	GenerateKey(hasher, method, u, headers, body)
	id := hex.EncodeToString(hasher.Sum(nil))
	return id[:10], filepath.Join(cacheBase, u.Scheme, u.Host, id)
}

func GenerateKey(w io.Writer, method string, u *url.URL, headers http.Header, body []byte) {
	if method == "" {
		method = http.MethodGet
	}
	w.Write([]byte(method))
	w.Write([]byte{'\r', '\n'})

	w.Write([]byte(u.Path))
	if u.RawQuery != "" {
		w.Write([]byte{'?'})
		w.Write([]byte(u.Query().Encode()))
	}
	w.Write([]byte{'\r', '\n'})

	headers.Write(w)
	w.Write([]byte{'\r', '\n'})

	w.Write(body)
	w.Write([]byte{'\r', '\n'})
}
