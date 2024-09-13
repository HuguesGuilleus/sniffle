package fetch

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Use only the cache
func CacheOnly(cacheBase string) Fetcher { return cacheOnly(cacheBase) }

type cacheOnly string

func (cacheOnly) Name() string { return "cache" }

func (cacheBase cacheOnly) Fetch(ctx context.Context, method string, u *url.URL, headers http.Header, body []byte) (io.ReadCloser, string, error) {
	logId, filePath := GeneratePath(string(cacheBase), method, u, headers, body)
	f, err := os.Open(filePath)
	return f, logId, err
}
