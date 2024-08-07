// Interface to write file
package writefile

import (
	"os"
	"path/filepath"
)

type WriteFile interface {
	// Path use '/' independant of os.
	// It's to implementation role to manage it.
	WriteFile(path string, data []byte) error
}

func Os(base string) WriteFile { return osBase(base) }

type osBase string

func (base osBase) WriteFile(path string, data []byte) error {
	path = filepath.Join(string(base), filepath.FromSlash(path))

	if err := os.MkdirAll(filepath.Dir(path), 0o775); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o664)
}

// A test map to check if a file (no empty) is written
type T map[string]bool

// Set t[path]=true if data is no empty.
// Alway return nil.
func (t T) WriteFile(path string, data []byte) error {
	if len(data) > 0 {
		t[path] = true
	}
	return nil
}
