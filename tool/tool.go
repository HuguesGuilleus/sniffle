package tool

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/url"
	"os"
	"sniffle/tool/fetch"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"sniffle/tool/writefile"
	"strings"
	"sync"
	"sync/atomic"
)

type Config struct {
	Logger *slog.Logger

	HostURL   string
	Languages []language.Language

	Writefile writefile.WriteFile
	Fetcher   []fetch.Fetcher

	Dev bool
}

func New(config *Config) *Tool {
	return &Tool{
		Logger: config.Logger,

		HostURL:   strings.TrimRight(config.HostURL, "/"),
		Languages: config.Languages,

		langRedirect: langRedirect(config.Languages),

		writefile: config.Writefile,
		fetcher:   config.Fetcher,

		dev: config.Dev,
	}
}

// All information for one service.
//
// Concerning WriteFile and Fetch:
// - log the result
// - consider error as fatal: return nothing because the service cannot resolve it.
type Tool struct {
	*slog.Logger

	// The host URL, whithout the training slash.
	// Ex: https://www.example.com
	HostURL   string
	Languages []language.Language

	// A file, to be put as .../index.html to redirect user with js to .../[lang].html page.
	// It generated from Languages.
	langRedirect []byte

	writeSum  uint64
	writefile writefile.WriteFile
	fetcher   []fetch.Fetcher

	// List of html files
	// Do not include index.html, because it's a js redirect
	// Do not include render.Back
	htmlFiles     []string
	htmlFileMutex sync.Mutex

	dev bool
}

// Is in dev mode or not?
func (t *Tool) Dev() bool { return t.dev }

func (t *Tool) WriteFile(path string, data []byte) {
	if strings.HasSuffix(path, ".html") && !(bytes.Equal(data, render.Back) || bytes.Equal(data, t.langRedirect)) {
		t.htmlFileMutex.Lock()
		t.htmlFiles = append(t.htmlFiles, path)
		t.htmlFileMutex.Unlock()
	}

	atomic.AddUint64(&t.writeSum, uint64(len(data)))

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

func NewTestTool(fetcher fetch.TestFetcher) (writefile.T, *Tool) {
	wf := writefile.T{}
	return wf, New(&Config{
		HostURL:   "https://example.com",
		Languages: []language.Language{language.English, language.French},

		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})),

		Writefile: wf,
		Fetcher:   []fetch.Fetcher{fetcher},
	})
}

// Write a ".../index.html" file that redirect user with js.
// Redirect will be ".../[lang].html".
func (t *Tool) LangRedirect(path string) {
	t.WriteFile(path, t.langRedirect)
}

func langRedirect(langs []language.Language) (h []byte) {
	h = []byte(`<!DOCTYPE html>` +
		`<html>` +
		`<head>` +
		`<meta charset=utf-8>` +
		`<meta name=robots content=noindex>` +
		`</head>` +
		`<body>` +
		`<p>Choose a language:</p>`)

	for _, l := range langs {
		h = append(h, `<a href=`...)
		h = append(h, l.String()...)
		h = append(h, `.html>`...)
		h = append(h, l.Human()...)
		h = append(h, `</a><br>`...)
	}

	h = append(h, `<script>`...)
	h = append(h, langRedirectJS[0]...)
	for i, l := range langs {
		if i != 0 {
			h = append(h, '/')
		}
		h = append(h, l.String()...)
	}
	h = append(h, langRedirectJS[1]...)
	h = append(h, `</script>`...)

	return
}

var langRedirectJS = func() [2]string {
	s := strings.NewReplacer("\n", "", "\t", "").Replace(`
		l="LANG".split('/'),
			u=navigator.languages,
			i=0,
			l;

		for(;i<u.length;i++){
			l=u[i].replace(/-\w+/,"");
			if(l.indexOf(l)+1){location=l+".html";break}
		}

		`)

	before, after, _ := strings.Cut(s, "LANG")

	return [2]string{before, after}
}()
