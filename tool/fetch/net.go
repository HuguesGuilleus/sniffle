package fetch

import (
	"bytes"
	"net/http"
	"sniffle/tool/writefs"
	"strings"
	"sync"
	"time"
)

type netFetcher struct {
	roundTripper http.RoundTripper
	cache        writefs.CreateOpener
	// Sleep interval between each request
	delayDuration map[string]time.Duration
	delayChannels map[string]<-chan struct{}
	delayMutex    sync.Mutex
}

func Net(roundTripper http.RoundTripper, cache writefs.CreateOpener, delay map[string]time.Duration) Fetcher {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}
	limit := make(chan struct{}, 1)
	limit <- struct{}{}
	if delay == nil {
		delay = make(map[string]time.Duration, 1)
	}
	delay[""] = max(delay[""], 0)
	return &netFetcher{
		roundTripper:  roundTripper,
		cache:         cache,
		delayDuration: delay,
		delayChannels: make(map[string]<-chan struct{}),
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

	defer fetcher.wait(request.URL.Host)()

	httpResponse, err := fetcher.roundTripper.RoundTrip(httpRequest)
	if err != nil {
		return nil, err
	}

	response := &Response{
		Status: httpResponse.StatusCode,
		Header: httpResponse.Header,
		Body:   httpResponse.Body,
	}

	path := request.Path()
	if f, err := fetcher.cache.Create(path); err != nil {
		return nil, err
	} else if err := SaveHTTP(request, response, time.Now(), f); err != nil {
		return nil, err
	} else if err := f.Close(); err != nil {
		return nil, err
	}

	f, err := fetcher.cache.Open(path)
	if err != nil {
		return nil, err
	}
	return ReadResponse(f)
}

func (fetcher *netFetcher) wait(host string) func() {
	delay, ok := time.Duration(0), false
	for h := host; !ok; _, h, _ = strings.Cut(h, ".") {
		delay, ok = fetcher.delayDuration[h]
	}
	if delay == 0 {
		return func() {}
	}

	newChannel := make(chan struct{})

	fetcher.delayMutex.Lock()
	previousChannel := fetcher.delayChannels[host]
	fetcher.delayChannels[host] = newChannel
	fetcher.delayMutex.Unlock()

	if previousChannel != nil {
		<-previousChannel
	}

	return func() {
		time.AfterFunc(delay, func() { newChannel <- struct{}{} })
	}
}
