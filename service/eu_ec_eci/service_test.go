package eu_ec_eci

import (
	"context"
	"io"
	"log/slog"
	"sniffle/tool"
	"sniffle/tool/writefile"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	wf, to := tool.NewTestTool(fetcher)
	to.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	Do(context.Background(), to)
	assert.Equal(t, writefile.T{
		"/eu/ec/eci/index.html":        1,
		"/eu/ec/eci/fr.html":           1,
		"/eu/ec/eci/en.html":           1,
		"/eu/ec/eci/schema.html":       1,
		"/eu/ec/eci/2024/index.html":   1,
		"/eu/ec/eci/2024/9/index.html": 1,
		"/eu/ec/eci/2024/9/logo.png":   1,
		"/eu/ec/eci/2024/9/fr.html":    1,
		"/eu/ec/eci/2024/9/en.html":    1,
	}, wf)
}
