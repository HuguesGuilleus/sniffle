package component

import (
	_ "embed"
	"sniffle/front/translate"
	"sniffle/tool/fronttool"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"time"
)

const (
	JsSearch uint = 1 << iota
	JsToc
)

var (
	//go:embed js/0.js
	raw0 []byte
	//go:embed js/search.js
	rawSearch []byte
	//go:embed js/toc.js
	rawToc []byte

	scripts = [...]render.H{
		0:                fronttool.InlineJs(raw0),
		JsSearch:         fronttool.InlineJs(raw0, rawSearch),
		JsToc:            fronttool.InlineJs(raw0, rawToc),
		JsSearch | JsToc: fronttool.InlineJs(raw0, rawSearch, rawToc),
	}
)

// A footer node. It should be the last element in the page.
// It contain in the end, so the DOM is complete when it's executed.
func Footer(l language.Language, flag uint) render.Node {
	return render.N("footer",
		translate.AllTranslation[l].FooterBuild,
		time.Now(),
		render.H("<br>"),
		render.Na("a", "href", "/about/"+l.String()+".html").N(
			translate.AllTranslation[l].AboutTextLink),
		scripts[flag],
	)
}
