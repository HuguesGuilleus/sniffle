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

type WriteReadFile interface {
	WriteFile
	// Like fs.ReadFileFS.
	// But the output bytes can be shared because it will not be modified.
	ReadFile(name string) ([]byte, error)
}

func Os(base string) WriteReadFile { return osBase(base) }

type osBase string

func (base osBase) WriteFile(path string, data []byte) error {
	path = filepath.Join(string(base), filepath.FromSlash(path))

	if err := os.MkdirAll(filepath.Dir(path), 0o775); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o664)
}

func (base osBase) ReadFile(path string) ([]byte, error) {
	path = filepath.Join(string(base), filepath.FromSlash(path))
	return os.ReadFile(path)
}

// A test map to check if a file (no empty) is written
type T map[string]int

// Set t[path]=true if data is no empty.
// Alway return nil.
func (t T) WriteFile(path string, data []byte) error {
	t[path]++
	return nil
}

// Return true if the writeFile is not called (len == 0).
func (t T) NoCalled() bool {
	return len(t) == 0
}
