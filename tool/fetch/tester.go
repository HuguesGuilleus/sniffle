package fetch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
)

type TestFetcher map[string][]byte

func (TestFetcher) Name() string { return "test" }

func (tf TestFetcher) Fetch(_ context.Context, u *url.URL) (io.ReadCloser, string, error) {
	logId, _ := GetFileID("", u)
	data, ok := tf[u.String()]
	if !ok {
		return nil, logId, fmt.Errorf("not found url %q", u.String())
	}
	return io.NopCloser(bytes.NewReader(data)), logId, nil
}
