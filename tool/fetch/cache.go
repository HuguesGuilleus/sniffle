package fetch

import (
	"context"
	"io"
	"net/url"
	"os"
)

// Use only the cache
func CacheOnly(cacheBase string) Fetcher { return cacheOnly(cacheBase) }

type cacheOnly string

func (cacheOnly) Name() string { return "cache" }

func (cacheBase cacheOnly) Fetch(ctx context.Context, u *url.URL) (io.ReadCloser, string, error) {
	logId, fileID := GetFileID(string(cacheBase), u)
	f, err := os.Open(fileID)
	return f, logId, err
}
