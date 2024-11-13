package fetch

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Clear cache.
//
// Each filter return a duration to keep the request.
// Is meta.Time < now+maxAge => remove it.
// The filter should the edit meta.
func ClearCache(cacheBase string, filters ...func(*Meta) time.Duration) error {
	m, err := filepath.Glob(filepath.Join(cacheBase, filepath.FromSlash("/*/*/*.http")))
	if err != nil {
		return err
	}

	now := time.Now()
	for _, p := range m {
		f, err := os.Open(p)
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
			if err := os.Remove(p); err != nil {
				return fmt.Errorf("remove cache %q: %w", p, err)
			}
		}
	}

	return nil
}
