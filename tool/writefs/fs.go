package writefs

import (
	"io"
	"io/fs"
)

type stdFS struct {
	fs fs.FS
}

// Transform a [fs.FS] from std lib to a [Opener] for this package.
func FS(fs fs.FS) Opener { return stdFS{fs} }

func (fsys stdFS) Open(path string) (io.ReadCloser, error) {
	return fsys.fs.Open(path)
}
