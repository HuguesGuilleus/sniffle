package eu_ec_eci

import (
	"io"
	"log/slog"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/language"
	"sniffle/tool/writefile"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	defer func(langs []language.Language) { translate.Langs = langs }(translate.Langs)
	translate.Langs = []language.Language{language.English, language.French}

	wf, to := tool.NewTestTool(fetcher)
	to.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	Do(to)
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
