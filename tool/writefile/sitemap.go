package writefile

import (
	"bytes"
	"slices"
	"sniffle/tool/render"
	"strings"
	"sync"
)

type SitemapWatcher struct {
	Sub WriteFile

	// Add is called on each .WriteFile, return true to keep the path.
	Add func(path string, data []byte) (keep bool)

	mutex sync.Mutex
	paths []string
}

// Sitemap replaces sub with a new SitemapWatcher and return it.
func Sitemap(sub *WriteFile) *SitemapWatcher {
	w := &SitemapWatcher{
		Sub: *sub,
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

func (w *SitemapWatcher) WriteFile(path string, data []byte) error {
	if w.Add(path, data) {
		w.mutex.Lock()
		w.paths = append(w.paths, path)
		w.mutex.Unlock()
	}

	return w.Sub.WriteFile(path, data)
}

// Sitemap generates a simple sitemap file (one url by line).
func (w *SitemapWatcher) Sitemap(host string) []byte {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	slices.Sort(w.paths)

	cap := 0
	for _, p := range w.paths {
		cap += len(host) + len(p) + 1
	}

	data := make([]byte, 0, cap)
	for _, p := range w.paths {
		data = append(data, host...)
		data = append(data, p...)
		data = append(data, '\n')
	}

	return data
}
