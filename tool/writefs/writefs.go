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

func WriteFile(c Creator, path string, data []byte) error {
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
