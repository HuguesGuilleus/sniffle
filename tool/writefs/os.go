package writefs

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type osFS string

func Os(dir string) CompleteFS { return osFS(filepath.Clean(dir)) }

func (o osFS) path(path string) string {
	return filepath.Join(string(o), filepath.FromSlash(path))
}

func (o osFS) Open(path string) (io.ReadCloser, error) {
	return os.Open(o.path(path))
}

func (o osFS) Create(path string) (io.WriteCloser, error) {
	p := o.path(path)
	if err := os.MkdirAll(filepath.Dir(p), 0o775); err != nil {
		return nil, err
	}
	return os.Create(p)
}

func (o osFS) FileIndex() ([]string, error) {
	list := make([]string, 0)
	err := fs.WalkDir(os.DirFS(string(o)), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			list = append(list, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (o osFS) Remove(path string) error {
	return os.Remove(o.path(path))
}
