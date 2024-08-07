package fetch

import (
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
	l := make(chan struct{}, max(limit, 1))
	for range cap(l) {
		l <- struct{}{}
	}

	return &netFetcher{roundTripper, cacheBase, l, sleep}
}

func (netFetcher) Name() string { return "net" }

func (fetcher netFetcher) Fetch(ctx context.Context, u *url.URL) (io.ReadCloser, string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "", nil)
	if err != nil {
		return nil, "", fmt.Errorf("make request fail: %w", err)
	}
	request.URL = u

	logId, fileID := GetFileID(fetcher.cacheBase, request.URL)
	if err := os.MkdirAll(filepath.Dir(fileID), 0o775); err != nil {
		return nil, logId, err
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
	saveMeta(fileID, &meta{u.String(), response.Status, time.Now(), response.Header})

	if response.StatusCode/100 != 2 {
		return nil, logId, fmt.Errorf("wrong status: %q", response.Status)
	}

	f, err := os.OpenFile(fileID, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o664)
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
	ResponseHeader http.Header `json:"responseHeader"`
}

func saveMeta(fileID string, m *meta) {
	j, _ := json.MarshalIndent(m, "", "\t")
	os.WriteFile(fileID+".json", j, 0o664)
}
