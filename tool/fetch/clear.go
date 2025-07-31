package fetch

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/HuguesGuilleus/sniffle/tool/writefs"
)

// Clear cache.
//
// Each filter return a duration to keep the request.
// If meta.Time < now+maxAge => remove it.
// The filter should not the edit meta.
func ClearCache(fsys writefs.CompleteFS, filters ...func(*Meta) time.Duration) error {
	paths, err := indexHTTPFiles(fsys)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, p := range paths {
		f, err := fsys.Open(p)
		if err != nil {
			return err
		}
		meta, err := ReadMeta(f)
		if err != nil {
			return fmt.Errorf("read meta cache %q: %w", p, err)
		}

		maxAge := time.Duration(0)
		for _, f := range filters {
			maxAge = max(maxAge, f(meta))
		}

		if meta.Time.Add(maxAge).Before(now) {
			if err := fsys.Remove(p); err != nil {
				return fmt.Errorf("remove cache %q: %w", p, err)
			}
		}
	}

	return nil
}

func indexHTTPFiles(fsys writefs.CompleteFS) ([]string, error) {
	index, err := fsys.FileIndex()
	if err != nil {
		return nil, err
	}
	index = slices.DeleteFunc(index, func(path string) bool { return !strings.HasSuffix(path, ".http") })
	slices.Sort(index)
	return index, nil
}
