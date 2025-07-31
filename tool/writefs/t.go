package writefs

import "io"

// A test [Creator] file system to check which file are created.
type T map[string]int

func (t T) Create(path string) (io.WriteCloser, error) {
	t[path]++
	return tw{}, nil
}

type tw struct{}

func (tw) Write(p []byte) (int, error) { return len(p), nil }
func (tw) Close() error                { return nil }
