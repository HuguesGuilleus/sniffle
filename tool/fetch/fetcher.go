package fetch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/url"
	"path/filepath"
)

type Fetcher interface {
	// Make a HTTP GET request and return body with cache ID.
	// If response code is out of [200,299] return an error.
	Fetch(ctx context.Context, u *url.URL) (io.ReadCloser, string, error)
	// Static name of the fetcher: net, cache...
	// Use this for the log
	Name() string
}

func GetFileID(cacheBase string, u *url.URL) (logID, fileID string) {
	query := ""
	if u.RawQuery != "" {
		query = "?" + u.Query().Encode()
	}
	h := sha256.Sum256([]byte(u.Path + query))
	id := hex.EncodeToString(h[:])
	return id[:10], filepath.Join(cacheBase, u.Scheme, u.Host, id)
}
