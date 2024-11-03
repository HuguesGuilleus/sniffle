package tool

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
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

	LongTasksCache writefile.WriteReadFile
	LongTasksMap   map[string]func([]byte) ([]byte, error)
}

func New(config *Config) *Tool {
	return &Tool{
		Logger: config.Logger,

		HostURL:   strings.TrimRight(config.HostURL, "/"),
		Languages: config.Languages,

		langRedirect: langRedirect(config.Languages),

		writefile: config.Writefile,
		fetcher:   config.Fetcher,

		longTasksCache: config.LongTasksCache,
		longTasksMap:   config.LongTasksMap,
	}
}

// All information for one service.
//
// Concerning WriteFile and Fetch:
// - log the result
// - consider error as fatal: return nothing because the service cannot resolve it.
type Tool struct {
	*slog.Logger

	// The host URL, without the training slash.
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

	longTasksCache writefile.WriteReadFile
	longTasksMutex sync.Mutex
	longTasksMap   map[string]func([]byte) ([]byte, error)
}

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
func (t *Tool) Fetch(ctx context.Context, method, urlString string, headers http.Header, body []byte) io.ReadCloser {
	u, err := url.Parse(urlString)
	if err != nil {
		t.Warn("http.failParseURL", "url", urlString, "err", err.Error())
		return nil
	}

	for _, f := range t.fetcher {
		r, id, err := f.Fetch(ctx, method, u, headers, body)
		if err != nil {
			t.Debug("http.fail", "fetcher", f.Name(), "id", id, "url", urlString, "err", err.Error())
			continue
		}
		t.Info("http.ok", "fetcher", f.Name(), "id", id, "url", urlString)
		return r
	}

	t.Warn("http.fatal", "url", urlString)
	return nil
}

// Try to execute a long task.
// The result can be empty if error occure or no task function.
func (t *Tool) LongTask(name, logRef string, input []byte) []byte {
	if t.longTasksCache == nil || t.longTasksMap == nil {
		t.Info("longtask.ignore", "name", name, "ref", logRef)
		return nil
	}

	idH := sha256.Sum256(input)
	id := hex.EncodeToString(idH[:])
	path := name + "/" + id

	if out, err := t.longTasksCache.ReadFile(path); err == nil {
		t.Info("longtask.cache", "name", name, "id", id, "ref", logRef)
		return out
	} else if !errors.Is(err, fs.ErrNotExist) {
		t.Warn("longtask.readErr", "name", name, "id", id, "ref", logRef, "err", err.Error())
	}

	t.longTasksMutex.Lock()
	task := t.longTasksMap[name]
	t.longTasksMutex.Unlock()
	if task == nil {
		t.Info("longtask.noFunc", "name", name, "id", id, "ref", logRef)
		return nil
	}

	out, err := task(input)
	if err != nil {
		t.Warn("longtask.err", "name", name, "id", id, "ref", logRef, "err", err.Error())
		return nil
	} else if err := t.longTasksCache.WriteFile(path, out); err != nil {
		t.Warn("longtask.writeErr", "name", name, "id", id, "ref", logRef, "err", err.Error())
	}

	t.Info("longtask.ok", "name", name, "id", id, "ref", logRef)
	return out
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

func langRedirect(langs []language.Language) []byte {
	langsStrings := make([]string, len(langs))
	for i, l := range langs {
		langsStrings[i] = l.String()
	}

	s := `<!DOCTYPE html>` +
		`<html>` +
		`<head>` +
		`<meta charset=utf-8>` +
		`<meta name=robots content=noindex>` +
		`<script>` +
		`for(a of navigator.languages)` +
		`if("` + strings.Join(langsStrings, "/") + `".split("/").includes(a=a.replace(/-\w+/,"")))` +
		`{location=a+".html";break}` +
		`</script>` +
		`</head>` +
		`<body>` +
		`<p>Choose a language:</p>`

	for _, l := range langs {
		s += `<a href=` + l.String() + `.html>` + l.Human() + `</a><br>`
	}

	return []byte(s)
}
