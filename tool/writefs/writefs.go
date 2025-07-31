// FS interface to read and write.
package writefs

import (
	"io"
)

type CreateOpener interface {
	Creator
	Opener
}

// A FS to open file and read the content.
type Opener interface {
	Open(path string) (io.ReadCloser, error)
}

func ReadAll(o Opener, path string) ([]byte, error) {
	r, err := o.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

// A FS to create and write file content.
type Creator interface {
	Create(path string) (io.WriteCloser, error)
}

// A sub type of [Creator] with one write.
type WriteFileFS interface {
	WriteFile(path string, data []byte) error
}

func WriteFile(c Creator, path string, data []byte) error {
	if fs, ok := c.(WriteFileFS); ok {
		return fs.WriteFile(path, data)
	}

	w, err := c.Create(path)
	if err != nil {
		return err
	}

	if _, err1 := w.Write(data); err1 != nil {
		w.Close()
		return err1
	}

	return w.Close()
}

// A full feature FS.
type CompleteFS interface {
	CreateOpener
	// Read files tree and return the index of all regular file.
	FileIndex() ([]string, error)
	// Remove path of a file.
	Remove(path string) error
}
