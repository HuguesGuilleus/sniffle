package front

import (
	_ "embed"
	"sniffle/front/frontcss"
	"sniffle/tool"
)

//go:embed favicon.ico
var favicon []byte

func WriteAssets(t *tool.Tool) {
	t.WriteFile("favicon.ico", favicon)
	t.WriteFile("style.css", frontcss.Style)
}
