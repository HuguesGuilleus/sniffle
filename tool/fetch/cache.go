package fetch

import (
	"github.com/HuguesGuilleus/sniffle/tool/writefs"
)

type cacheFetcher struct {
	fs writefs.Opener
}

func Cache(fs writefs.Opener) Fetcher { return &cacheFetcher{fs} }

func (*cacheFetcher) Name() string { return "cache" }

func (c *cacheFetcher) Fetch(request *Request) (*Response, error) {
	f, err := c.fs.Open(request.Path())
	if err != nil {
		return nil, err
	}
	return ReadResponse(f)
}
