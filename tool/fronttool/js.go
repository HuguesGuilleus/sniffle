package fronttool

import (
	"bytes"
	"sniffle/tool/render"

	"github.com/tdewolff/minify/v2/js"
)

// Create a inline JS from multiple files.
// Use strict JS and wrap all in a closure with document as argument.
// It panic if an minify error occurs.
func InlineJs(chuncks ...[]byte) render.H {
	src := bytes.Buffer{}
	src.WriteString(`"use strict";((document) => {`)
	for _, chunck := range chuncks {
		src.Write(chunck)
		src.WriteString("\n\n")
	}
	src.WriteString(`})(document);`)

	out := bytes.Buffer{}
	out.WriteString(`<script>`)
	if err := js.Minify(nil, &out, &src, nil); err != nil {
		panic(err)
	}
	out.WriteString(`</script>`)

	return render.H(out.String())
}
