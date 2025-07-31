package component

import (
	"testing"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/stretchr/testify/assert"
)

func TestLangAlternate(t *testing.T) {
	nodes := LangAlternate("/yolo/", language.French,
		[]language.Language{language.English, language.French})

	assert.Equal(t, `<!DOCTYPE html><_>`+
		`<link rel=alternate hreflang=en href=https://sniffle.eu/yolo/en.html>`+
		`</_>`,
		string(render.Merge(render.N("_", nodes))))
}
