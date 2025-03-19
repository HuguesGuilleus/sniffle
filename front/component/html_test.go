package component

import (
	"sniffle/common/language"
	"sniffle/tool/render"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLangAlternate(t *testing.T) {
	nodes := LangAlternate("https://sniffle.eu/yolo/", language.French,
		[]language.Language{language.English, language.French})

	assert.Equal(t, `<!DOCTYPE html><_>`+
		`<link rel=alternate hreflang=en href=https://sniffle.eu/yolo/en.html>`+
		`</_>`,
		string(render.Merge(render.N("_", nodes))))
}
