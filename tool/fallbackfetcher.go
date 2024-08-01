package tool

import (
	"context"
	"log/slog"
	"net/http"
)

// Create a fetcher, use cache and if error, use a standard fetcher.
func FallBackFetcher(logger *slog.Logger, roundTripper http.RoundTripper, cacheBase string, limit int) Fetcher {
	return &fallBackFetcher{
		read: readFetcher{logger, cacheBase},
		std:  NewFetcher(logger, roundTripper, cacheBase, limit),
	}
}

type fallBackFetcher struct {
	read readFetcher
	std  Fetcher
}

func (f *fallBackFetcher) FetchGET(ctx context.Context, u string) ([]byte, error) {
	data, err := f.read.FetchGET(ctx, u)
	if err == nil {
		return data, nil
	}

	return f.std.FetchGET(ctx, u)
}

func (f *fallBackFetcher) Error(msg string, args ...any) {
	f.read.logger.Error(msg, args...)
}
