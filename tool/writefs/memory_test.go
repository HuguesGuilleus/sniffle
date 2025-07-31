package writefs

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	m := Memory()

	// WriteFile
	assert.NoError(t, m.(WriteFileFS).WriteFile("/file.txt", []byte("hello")))
	data, err := ReadAll(m, "/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(data))

	// Create
	w, err := m.Create("/w.txt")
	assert.NoError(t, err)
	_, err = w.Write([]byte("Hello World!"))
	assert.NoError(t, err)
	assert.NoError(t, w.Close())

	data, err = ReadAll(m, "/file.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(data))

	// Remove
	assert.NoError(t, m.Remove("/file.txt"))

	// Open not existing file
	data, err = ReadAll(m, "/file.txt")
	assert.ErrorIs(t, err, fs.ErrNotExist)
	assert.Nil(t, data)

	// Index
	index, err := m.FileIndex()
	assert.NoError(t, err)
	assert.Equal(t, []string{"/w.txt"}, index)
}
