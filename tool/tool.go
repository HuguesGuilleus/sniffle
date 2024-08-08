package tool

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/url"
	"sniffle/tool/fetch"
	"sniffle/tool/language"
	"sniffle/tool/writefile"
)

type Config struct {
	Logger *slog.Logger

	HostURL   string
	Languages []language.Language

	Writefile writefile.WriteFile
	Fetcher   []fetch.Fetcher
}

func New(config *Config) *Tool {
	return &Tool{
		Logger: config.Logger,

		HostURL:   config.HostURL,
		Languages: config.Languages,

		writefile: config.Writefile,
		fetcher:   config.Fetcher,
	}
}

// All information for one service.
//
// Concerning WriteFile and Fetch:
// - log the result
// - consider error as fatal: return nothing because the service cannot resolve it.
type Tool struct {
	*slog.Logger

	HostURL   string
	Languages []language.Language

	writefile writefile.WriteFile
	fetcher   []fetch.Fetcher
}

func (t *Tool) WriteFile(path string, data []byte) {
	err := t.writefile.WriteFile(path, data)
	if err != nil {
		t.Warn("out.fail", "path", path, "err", err.Error())
	} else {
		t.Debug("out.ok", "path", path, "len", len(data))
	}
}

// Make a HTTP call, and return the body.
// Return nil if cannot parse the url.
//
// Logs results and errors.
//
// If all internal fetchers return and error, return nil.
func (t *Tool) Fetch(ctx context.Context, urlString string) io.ReadCloser {
	u, err := url.Parse(urlString)
	if err != nil {
		t.Warn("http.failParseURL", "url", urlString, "err", err.Error())
		return nil
	}

	for _, f := range t.fetcher {
		r, id, err := f.Fetch(ctx, u)
		if err != nil {
			t.Debug("http.fail", "fetcher", f.Name(), "id", id, "url", urlString, "err", err.Error())
			continue
		}
		t.Info("http.ok", "fetcher", f.Name(), "id", id, "url", urlString)
		return r
	}

	t.Debug("http.fatal", "url", urlString)
	return nil
}

// For dev transistion
func (t *Tool) FetchGET(ctx context.Context, urlString string) ([]byte, error) {
	reader := t.Fetch(ctx, urlString)
	if reader == nil {
		return nil, errors.New("fail")
	}
	return io.ReadAll(reader)
}
