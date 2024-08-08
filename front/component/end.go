package component

import (
	"bytes"
	_ "embed"
	"sniffle/tool/render"

	"github.com/tdewolff/minify/v2/js"
)

var (
	//go:embed end.js
	end []byte
	End render.H = func() render.H {
		buff := bytes.Buffer{}
		buff.WriteString(`<script defer>`)
		if err := js.Minify(nil, &buff, bytes.NewReader(end), nil); err != nil {
			panic(err)
		}
		buff.WriteString(`</script>`)
		return render.H(buff.String())
	}()
)
