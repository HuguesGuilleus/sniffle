package tool

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"sniffle/tool/fetch"
	"sniffle/tool/writefs"
	"strings"
	"time"
)

const NoticeLevel = slog.LevelInfo + 2

type Config struct {
	LogHandler slog.Handler

	Writefile writefs.Creator
	Fetcher   []fetch.Fetcher

	LongTasksCache writefs.CreateOpener
	LongTasksMap   map[string]func(*Tool, []byte) ([]byte, error)
}

func (config *Config) Run(name string, do func(*Tool)) {
	if strings.HasPrefix(name, "//") && !DevMode {
		return
	}

	begin := time.Now()

	t := New(config)
	t.Logger = t.Logger.With(slog.Any("service", name))

	do(t)

	t.Log(context.Background(), NoticeLevel, "end", "duration", time.Since(begin))
}

// All information for one service.
//
// Concerning WriteFile and Fetch:
// - log the result
// - consider error as fatal: return nothing because the service cannot resolve it.
type Tool struct {
	*slog.Logger

	writefile writefs.Creator
	fetcher   []fetch.Fetcher

	longTasksCache writefs.CreateOpener
	longTasksMap   map[string]func(*Tool, []byte) ([]byte, error)
}

func New(config *Config) *Tool {
	return &Tool{
		Logger: slog.New(config.LogHandler),

		writefile: config.Writefile,
		fetcher:   config.Fetcher,

		longTasksCache: config.LongTasksCache,
		longTasksMap:   config.LongTasksMap,
	}
}

func (t *Tool) WriteFile(path string, data []byte) {
	err := writefs.WriteFile(t.writefile, path, data)
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

	if out, err := writefs.ReadAll(t.longTasksCache, path); err == nil {
		t.Info("longtask.cache", "name", name, "id", id, "ref", logRef)
		return out
	} else if !errors.Is(err, fs.ErrNotExist) {
		t.Warn("longtask.cacheErr", "name", name, "id", id, "ref", logRef, "err", err.Error())
	}

	task := t.longTasksMap[name]
	if task == nil {
		t.Info("longtask.noFunc", "name", name, "id", id, "ref", logRef)
		return nil
	}

	out, err := task(t, input)
	if err != nil {
		t.Warn("longtask.err", "name", name, "id", id, "ref", logRef, "err", err.Error())
		return nil
	} else if err := writefs.WriteFile(t.longTasksCache, path, out); err != nil {
		t.Warn("longtask.writeErr", "name", name, "id", id, "ref", logRef, "err", err.Error())
	}

	t.Info("longtask.ok", "name", name, "id", id, "ref", logRef)
	return out
}

func NewTestTool(fetcher fetch.Fetcher) (writefs.T, *Tool) {
	wfs := writefs.T{}
	return wfs, New(&Config{
		LogHandler: slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}),

		Writefile: wfs,
		Fetcher:   []fetch.Fetcher{fetcher},
	})
}
