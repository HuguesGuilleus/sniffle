package fetch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type netFetcher struct {
	http.RoundTripper
	cacheBase string
	limit     chan struct{}
	// Sleep interval between each request
	sleep time.Duration
}

func Net(roundTripper http.RoundTripper, cacheBase string, limit uint, sleep time.Duration) Fetcher {
	if roundTripper == nil {
		roundTripper = http.DefaultTransport
	}

	l := make(chan struct{}, max(limit, 1))
	for range cap(l) {
		l <- struct{}{}
	}

	return &netFetcher{roundTripper, cacheBase, l, sleep}
}

func (netFetcher) Name() string { return "net" }

func (fetcher netFetcher) Fetch(ctx context.Context, method string, u *url.URL, headers http.Header, body []byte) (io.ReadCloser, string, error) {
	logId, filePath := GeneratePath(fetcher.cacheBase, method, u, headers, body)
	if err := os.MkdirAll(filepath.Dir(filePath), 0o775); err != nil {
		return nil, logId, err
	}

	request, err := http.NewRequestWithContext(ctx, method, "", bytes.NewReader(body))
	if err != nil {
		return nil, logId, err
	}
	request.URL = u
	for k, vv := range headers {
		for _, v := range vv {
			request.Header.Add(k, v)
		}
	}

	<-fetcher.limit
	defer func() {
		time.Sleep(fetcher.sleep)
		fetcher.limit <- struct{}{}
	}()

	response, err := fetcher.RoundTrip(request)
	if err != nil {
		return nil, logId, err
	}
	defer response.Body.Close()
	saveMeta(filePath, &meta{u.String(), response.Status, time.Now(), string(body), request.Header, response.Header})

	if response.StatusCode/100 != 2 {
		return nil, logId, fmt.Errorf("wrong status: %q", response.Status)
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o664)
	if err != nil {
		return nil, logId, err
	} else if _, err := io.Copy(f, response.Body); err != nil {
		f.Close()
		return nil, logId, err
	} else if _, err := f.Seek(0, io.SeekStart); err != nil {
		f.Close()
		return nil, logId, err
	}

	return f, logId, nil
}

type meta struct {
	URL            string      `json:"url"`
	Status         string      `json:"status"`
	Time           time.Time   `json:"time"`
	Body           string      `json:"body"`
	RequestHeader  http.Header `json:"requestHeader"`
	ResponseHeader http.Header `json:"responseHeader"`
}

func saveMeta(filePath string, m *meta) {
	j, _ := json.MarshalIndent(m, "", "\t")
	os.WriteFile(filePath+".json", j, 0o664)
}
