package tool

import (
	"context"
	"encoding/json"
	"io"
)

func FetchJSON(ctx context.Context, t *Tool, url string, v any) (fail bool) {
	reader := t.Fetch(ctx, url)
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

func FetchAll(ctx context.Context, t *Tool, url string) []byte {
	reader := t.Fetch(ctx, url)
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
