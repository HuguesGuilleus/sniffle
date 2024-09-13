package tool

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func FetchJSON(ctx context.Context, t *Tool, method, url string, headers http.Header, body []byte, v any) (fail bool) {
	reader := t.Fetch(ctx, method, url, headers, body)
	if reader == nil {
		return true
	}
	defer reader.Close()

	if err := json.NewDecoder(reader).Decode(v); err != nil {
		t.Warn("http.decodeJsonFail", "url", url, "err", err.Error())
		return true
	}

	return false
}

func FetchAll(ctx context.Context, t *Tool, method, url string, headers http.Header, body []byte) []byte {
	reader := t.Fetch(ctx, method, url, headers, body)
	if reader == nil {
		return nil
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Warn("http.readAllfail", "url", url, "err", err.Error())
		return nil
	}

	return data
}
