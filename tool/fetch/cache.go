package fetch

import "os"

type cache string

func Cache(cacheBase string) Fetcher { return cache(cacheBase) }

func (cache) Name() string { return "cache" }

func (cache cache) Fetch(request *Request) (*Response, error) {
	f, err := os.Open(getPath(string(cache), request))
	if err != nil {
		return nil, err
	}
	return ReadResponse(f)
}
