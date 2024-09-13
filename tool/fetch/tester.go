package fetch

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TestFetcher map[string][]byte

func (TestFetcher) Name() string { return "test" }

func (tf TestFetcher) Fetch(_ context.Context, method string, u *url.URL, headers http.Header, body []byte) (io.ReadCloser, string, error) {
	logId, _ := GeneratePath("", method, u, headers, body)

	key := bytes.Buffer{}
	key.WriteString(u.String())
	if (method != "" && method != "GET") || len(headers) != 0 || len(body) != 0 {
		key.WriteByte('\n')
		GenerateKey(&key, method, u, headers, body)
	}

	data, ok := tf[key.String()]
	if !ok {
		return nil, logId, fmt.Errorf("not found url %q", key.String())
	}
	return io.NopCloser(bytes.NewReader(data)), logId, nil
}
