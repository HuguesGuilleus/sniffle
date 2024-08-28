package component

import (
	"bytes"
	_ "embed"
	"sniffle/front/translate"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"time"

	"github.com/tdewolff/minify/v2/js"
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
		0:                minifyJs(raw0),
		JsSearch:         minifyJs(raw0, rawSearch),
		JsToc:            minifyJs(raw0, rawToc),
		JsSearch | JsToc: minifyJs(raw0, rawSearch, rawToc),
	}
)

func minifyJs(chuncks ...[]byte) render.H {
	src := bytes.Buffer{}
	src.WriteString(`((document) => {`)
	for _, chunck := range chuncks {
		src.Write(chunck)
		src.WriteString("\n\n")
	}
	src.WriteString(`})(document);`)

	buff := bytes.Buffer{}
	buff.WriteString(`<script>"use strict";`)
	bytes.Join(chuncks, []byte("\n\n"))
	if err := js.Minify(nil, &buff, &src, nil); err != nil {
		panic(err)
	}
	buff.WriteString(`</script>`)
	return render.H(buff.String())
}

// A footer node. It should be the last element in the page.
// It contain in the end, so the DOM is complete when it's executed.
func Footer(l language.Language, flag uint) render.Node {
	return render.N("footer",
		translate.AllTranslation[l].FooterBuild,
		time.Now(),
		render.H("<br>"),
		render.No("a", render.A("href", "/about/"+l.String()+".html"),
			translate.AllTranslation[l].AboutTextLink),
		scripts[flag],
	)
}
