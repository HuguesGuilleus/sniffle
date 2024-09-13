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
	assert.NoError(t, os.WriteFile("cache/https/example.com/498f07ec77e6610d3d5d527a9dda3dbcf81ae934b1a937afa5117ed5a21542b6", []byte(`body`), 0o664))

	fetcher := CacheOnly("cache")
	assert.Equal(t, "cache", fetcher.Name())
	body, id, err := fetcher.Fetch(context.Background(), "", u, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, "498f07ec77", id)
	defer body.Close()

	data, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.EqualValues(t, "body", data)
	assert.IsType(t, &os.File{}, body)

	assert.NoError(t, os.RemoveAll("cache"))
}
