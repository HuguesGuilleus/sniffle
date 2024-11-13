package fetch_test

import (
	"io"
	"sniffle/tool/fetch"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTester(t *testing.T) {
	f := fetch.Test(map[string]*fetch.TestResponse{
		"https://example.org/dir/file": fetch.TRs(200, "Hello World!", "Content-Type", "text/plain; charset=utf-8"),

		"PATCH\nhttps://example.org/api/title\n" +
			"Content-Type: application/json\n" +
			"\n" +
			`{"title":"Yolo"}` + "\n": fetch.TR(205, nil, "X-Header"),
	})
	assert.Equal(t, "test", f.Name())

	response, err := f.Fetch(fetch.R("", "https://example.org/dir/file", nil))
	assert.NoError(t, err)
	assert.Equal(t, *fetch.TRs(200, "Hello World!", "Content-Type", "text/plain; charset=utf-8"), fromResponse(response))

	response, err = f.Fetch(fetch.Rs("PATCH", "https://example.org/api/title",
		`{"title":"Yolo"}`,
		"Content-Type", "application/json",
	))
	assert.NoError(t, err)
	assert.Equal(t, *fetch.TR(205, []byte{}), fromResponse(response))
}

func fromResponse(r *fetch.Response) fetch.TestResponse {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	return fetch.TestResponse{
		Status: r.Status,
		Header: r.Header,
		Body:   body,
	}
}
