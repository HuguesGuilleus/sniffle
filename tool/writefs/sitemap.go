package writefs

import (
	"bytes"
	"io"
	"slices"
	"strings"
	"sync"

	"github.com/HuguesGuilleus/sniffle/tool/render"
)

// SitemapWatcher collect all HTML path to create a sitemap.
type SitemapWatcher struct {
	Creator Creator

	// Add is called on each .WriteFile, return true to keep the path.
	Add func(path string, data []byte) (keep bool)

	mutex sync.Mutex
	paths []string
}

// Sitemap replaces sub with a new SitemapWatcher and return it.
func Sitemap(sub *Creator) *SitemapWatcher {
	w := &SitemapWatcher{
		Creator: *sub,
		Add: func(path string, data []byte) (keep bool) {
			return strings.HasSuffix(path, ".html") &&
				!bytes.Equal(data, render.Back) &&
				!bytes.Contains(data, []byte(`<meta name=robots content=noindex>`))
		},
		paths: make([]string, 0),
	}
	*sub = w
	return w
}

func (watcher *SitemapWatcher) Create(path string) (io.WriteCloser, error) {
	return &sitemapWriter{
		watcher: watcher,
		path:    path,
	}, nil
}

type sitemapWriter struct {
	watcher *SitemapWatcher
	path    string
	bytes.Buffer
}

func (w *sitemapWriter) Close() error {
	return w.watcher.WriteFile(w.path, w.Bytes())
}

func (watcher *SitemapWatcher) WriteFile(path string, data []byte) error {
	if watcher.Add(path, data) {
		watcher.mutex.Lock()
		watcher.paths = append(watcher.paths, path)
		watcher.mutex.Unlock()
	}
	return WriteFile(watcher.Creator, path, data)
}

// Sitemap generates a simple sitemap file (one url by line).
//
// The host prefix all line.
// Host must not end by a '/'. Valid example: `sitemap.Sitemap("https://example.com")`
func (watcher *SitemapWatcher) Sitemap(host string) []byte {
	watcher.mutex.Lock()
	defer watcher.mutex.Unlock()

	slices.Sort(watcher.paths)

	cap := 0
	for _, p := range watcher.paths {
		cap += len(host) + len(p) + 1
	}

	data := make([]byte, 0, cap)
	for _, p := range watcher.paths {
		data = append(data, host...)
		data = append(data, p...)
		data = append(data, '\n')
	}

	return data
}
