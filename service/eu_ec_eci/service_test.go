package eu_ec_eci

import (
	"io"
	"log/slog"
	"testing"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/writefs"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	defer func(langs []language.Language) { translate.Langs = langs }(translate.Langs)
	translate.Langs = []language.Language{language.English, language.French}

	wfs, to := tool.NewTestTool(fetcher)
	to.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	Do(to)
	assert.Equal(t, writefs.T{
		"/eu/ec/eci/index.html":                 1,
		"/eu/ec/eci/fr.html":                    1,
		"/eu/ec/eci/en.html":                    1,
		"/eu/ec/eci/schema.html":                1,
		"/eu/ec/eci/data-extradelay/index.html": 1,
		"/eu/ec/eci/data-extradelay/fr.html":    1,
		"/eu/ec/eci/data-extradelay/en.html":    1,
		"/eu/ec/eci/data-threshold/index.html":  1,
		"/eu/ec/eci/data-threshold/fr.html":     1,
		"/eu/ec/eci/data-threshold/en.html":     1,
		"/eu/ec/eci/2024/index.html":            1,
		"/eu/ec/eci/2024/9/index.html":          1,
		"/eu/ec/eci/2024/9/fr.html":             1,
		"/eu/ec/eci/2024/9/en.html":             1,
		"/eu/ec/eci/refused/index.html":         1,
		"/eu/ec/eci/refused/en.html":            1,
		"/eu/ec/eci/refused/fr.html":            1,
		"/eu/ec/eci/refused/42/index.html":      1,
		"/eu/ec/eci/refused/42/en.html":         1,
		"/eu/ec/eci/refused/42/fr.html":         1,
	}, wfs)
}
