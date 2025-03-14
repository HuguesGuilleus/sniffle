package tool

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"sniffle/tool/fetch"
	"sniffle/tool/render"
	"sniffle/tool/writefile"
	"strings"
	"sync"
	"sync/atomic"
)

type Config struct {
	Logger *slog.Logger

	HostURL string

	Writefile writefile.WriteFile
	Fetcher   []fetch.Fetcher

	LongTasksCache writefile.WriteReadFile
	LongTasksMap   map[string]func([]byte) ([]byte, error)
}

func New(config *Config) *Tool {
	return &Tool{
		Logger: config.Logger,

		HostURL: strings.TrimRight(config.HostURL, "/"),

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
	HostURL string

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
	if strings.HasSuffix(path, ".html") && !(bytes.Equal(data, render.Back) || bytes.Contains(data, []byte(`<meta name=robots content=noindex>`))) {
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
func (t *Tool) Fetch(request *fetch.Request) *fetch.Response {
	id := request.ID()
	u := request.URL.String()
	for _, f := range t.fetcher {
		r, err := f.Fetch(request)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				t.Debug("http.notExist", "fetcher", f.Name(), "id", id, "url", u)
			} else {
				t.Warn("http.fail", "fetcher", f.Name(), "id", id, "url", u, "err", err.Error())
			}
			continue
		}
		t.Info("http.ok", "fetcher", f.Name(), "id", id[:8], "url", u, "status", r.Status)
		return r
	}

	t.Warn("http.fatal", "url", u)
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
		t.Warn("longtask.cacheErr", "name", name, "id", id, "ref", logRef, "err", err.Error())
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

func NewTestTool(fetcherMap map[string]*fetch.TestResponse) (writefile.T, *Tool) {
	wf := writefile.T{}
	return wf, New(&Config{
		HostURL: "https://example.com",

		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})),

		Writefile: wf,
		Fetcher:   []fetch.Fetcher{fetch.Test(fetcherMap)},
	})
}
