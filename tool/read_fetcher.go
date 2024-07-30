package tool

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
)

type readFetcher struct {
	logger    *slog.Logger
	cacheBase string
}

func ReadFetcher(logger *slog.Logger, cacheBase string) Fetcher {
	return readFetcher{logger, cacheBase}
}

func (r readFetcher) FetchGET(ctx context.Context, u string) ([]byte, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		r.logger.Error("http.wrongURL", "u", u, "err", err)
		return nil, fmt.Errorf("wrong URL syntax: %w", err)
	}
	p := filepath.Join(r.cacheBase, parsedURL.Scheme, parsedURL.Host, hashURL(parsedURL))

	data, err := os.ReadFile(p)
	if err != nil {
		r.logger.Error("http.readFile", "url", u, "err", err)
		return nil, fmt.Errorf("read for url %q fail: %w", u, err)
	}

	r.logger.Debug("http.ok", "url", u)
	return data, nil
}

func (r readFetcher) Error(msg string, args ...any) {
	r.logger.Error(msg, args...)
}
