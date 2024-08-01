package tool

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Fetcher interface {
	FetchGET(ctx context.Context, url string) ([]byte, error)
	Error(msg string, args ...any)
}

// Create a new standard fetcher.
// - Log each HTTP request and error
// - Put reponse in cacheBase
// - Do not exed limit argument of concurent request
func NewFetcher(logger *slog.Logger, roundTripper http.RoundTripper, cacheBase string, limit int) Fetcher {
	f := &stdFetcher{logger, roundTripper, cacheBase, make(chan struct{}, limit)}

	for range limit {
		f.limit <- struct{}{}
	}

	return f
}

// Standard fetch
// - log
// - put to cache (but do not read)
type stdFetcher struct {
	logger       *slog.Logger
	roundTripper http.RoundTripper
	cacheBase    string
	limit        chan struct{}
}

func (fetcher *stdFetcher) FetchGET(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fetcher.logger.Error("http.makeRequest", "url", url, "err", err.Error())
		return nil, err
	}

	<-fetcher.limit
	defer func() {
		fetcher.limit <- struct{}{}
	}()

	response, err := fetcher.roundTripper.RoundTrip(request)
	if err != nil {
		fetcher.logger.Error("http.send", "url", url, "err", err.Error())
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		fetcher.logger.Error("http.status", "url", url, "status", response.Status)
		return nil, fmt.Errorf("wrong HTTP status: %q", response.Status)
	}

	// Read body
	buff := bytes.Buffer{}
	if response.ContentLength > 0 {
		buff.Grow(int(response.ContentLength))
	}
	if _, err := buff.ReadFrom(response.Body); err != nil {
		fetcher.logger.Error("http.readBody", "url", url, "err", err.Error())
		return nil, err
	}

	// Save cache
	pathHex := hashURL(request.URL)
	cacheDir := filepath.Join(fetcher.cacheBase, request.URL.Scheme, request.URL.Host)
	cachePath := filepath.Join(fetcher.cacheBase, request.URL.Scheme, request.URL.Host, pathHex)

	meta, _ := json.Marshal(struct {
		URL            string      `json:"url"`
		Status         string      `json:"status"`
		Time           time.Time   `json:"time"`
		ResponseHeader http.Header `json:"responseHeader"`
	}{url, response.Status, time.Now(), response.Header})

	if err := os.MkdirAll(cacheDir, 0o775); err != nil {
		fetcher.logger.Error("http.saveCacheDir", "path", cacheDir, "err", err.Error())
		return nil, err
	} else if err := os.WriteFile(cachePath+".json", meta, 0o664); err != nil {
		fetcher.logger.Error("http.saveCacheMeta", "path", cachePath+".json", "err", err.Error())
		return nil, err
	} else if err := os.WriteFile(cachePath, buff.Bytes(), 0o664); err != nil {
		fetcher.logger.Error("http.saveCacheBlod", "path", cachePath, "err", err.Error())
		return nil, err
	}

	fetcher.logger.Info("http.ok", "url", url, "h", pathHex)

	return buff.Bytes(), err
}

func (fetcher *stdFetcher) Error(msg string, args ...any) {
	fetcher.logger.Error(msg, args...)
}

func hashURL(u *url.URL) string {
	query := ""
	if u.RawQuery != "" {
		query = "?" + u.Query().Encode()
	}
	hashArray := sha256.Sum256([]byte(u.Path + query))
	return hex.EncodeToString(hashArray[:])
}

func FetchGETJSON(ctx context.Context, fetcher Fetcher, url string, v any) error {
	data, err := fetcher.FetchGET(ctx, url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		fetcher.Error("http.parseJson", "url", url, "err", err.Error())
		return err
	}

	return nil
}
