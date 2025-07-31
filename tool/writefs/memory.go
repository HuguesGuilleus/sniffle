package writefs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"slices"
	"sync"
)

type memoryFS struct {
	files map[string][]byte
	mutex sync.RWMutex
}

func Memory() CompleteFS { return &memoryFS{files: make(map[string][]byte)} }

type memoryWriter struct {
	fs   *memoryFS
	path string
	bytes.Buffer
}

func (mw *memoryWriter) Close() error {
	mw.fs.mutex.Lock()
	defer mw.fs.mutex.Unlock()
	mw.fs.files[mw.path] = mw.Bytes()
	return nil
}

func (mem *memoryFS) Create(path string) (io.WriteCloser, error) {
	return &memoryWriter{
		fs:   mem,
		path: path,
	}, nil
}

func (mem *memoryFS) WriteFile(path string, data []byte) error {
	data = bytes.Clone(data)
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	mem.files[path] = data
	return nil
}

func (mem *memoryFS) Open(path string) (io.ReadCloser, error) {
	mem.mutex.RLock()
	defer mem.mutex.RUnlock()
	data, ok := mem.files[path]
	if !ok {
		return nil, fmt.Errorf("file %q: %w", path, fs.ErrNotExist)
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

func (mem *memoryFS) FileIndex() ([]string, error) {
	mem.mutex.RLock()
	defer mem.mutex.RUnlock()
	index := slices.AppendSeq(make([]string, 0, len(mem.files)), maps.Keys(mem.files))
	return index, nil
}

func (mem *memoryFS) Remove(path string) error {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()
	_, exist := mem.files[path]
	if !exist {
		return fmt.Errorf("file %q: %w", path, fs.ErrNotExist)
	}
	delete(mem.files, path)
	return nil
}
