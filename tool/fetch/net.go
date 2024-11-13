package fetch

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type netFetcher struct {
	roundTripper http.RoundTripper
	cacheBase    string
	limit        chan struct{}
	// Sleep interval between each request
	sleep time.Duration
}

func Net(roundTripper http.RoundTripper, cacheBase string, sleep time.Duration) Fetcher {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}
	limit := make(chan struct{}, 1)
	limit <- struct{}{}
	return &netFetcher{
		roundTripper: roundTripper,
		cacheBase:    filepath.Clean(cacheBase),
		limit:        limit,
		sleep:        sleep,
	}
}

func (*netFetcher) Name() string { return "net" }

func (fetcher *netFetcher) Fetch(request *Request) (*Response, error) {
	httpRequest, err := http.NewRequest(request.Method, "", bytes.NewReader(request.Body))
	if err != nil {
		return nil, err
	}
	httpRequest.URL = request.URL
	httpRequest.Header = request.Header.Clone()

	<-fetcher.limit
	defer func() {
		go func() {
			time.Sleep(fetcher.sleep)
			fetcher.limit <- struct{}{}
		}()
	}()

	httpResponse, err := fetcher.roundTripper.RoundTrip(httpRequest)
	if err != nil {
		return nil, err
	}

	response := &Response{
		Status: httpResponse.StatusCode,
		Header: httpResponse.Header,
		Body:   httpResponse.Body,
	}

	if err := os.MkdirAll(getDir(fetcher.cacheBase, request), 0o775); err != nil {
		return nil, err
	}
	path := getPath(fetcher.cacheBase, request)
	if f, err := os.Create(path); err != nil {
		return nil, err
	} else if err := SaveHTTP(request, response, time.Now(), f); err != nil {
		return nil, err
	} else if err := f.Close(); err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ReadResponse(f)
}
