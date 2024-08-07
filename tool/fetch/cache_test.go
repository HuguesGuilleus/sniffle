package fetch

import (
	"context"
	"io"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	u, err := url.Parse("https://example.com/file?b=2&a=1")
	assert.NoError(t, err)

	assert.NoError(t, os.MkdirAll("cache/https/example.com/ ", 0o775))
	assert.NoError(t, os.WriteFile("cache/https/example.com/62caa696594967cb5cd8957da8eb44fbced0b6b9acd0b2ad306ab70b05385fd5", []byte(`body`), 0o664))

	fetcher := CacheOnly("cache")
	assert.Equal(t, "cache", fetcher.Name())
	body, id, err := fetcher.Fetch(context.Background(), u)
	assert.NoError(t, err)
	assert.Equal(t, "62caa69659", id)
	defer body.Close()

	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.EqualValues(t, "body", data)
	assert.IsType(t, &os.File{}, body)

	assert.NoError(t, os.RemoveAll("cache"))
}
