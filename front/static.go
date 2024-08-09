package front

import (
	"context"
	_ "embed"
	"sniffle/front/frontcss"
	"sniffle/tool"
)

//go:embed favicon.ico
var favicon []byte

//go:embed robots.txt
var robots []byte

func Do(_ context.Context, t *tool.Tool) {
	t.WriteFile("favicon.ico", favicon)
	t.WriteFile("robots.txt", robots)
	t.WriteFile("style."+frontcss.StyleHash+".css", frontcss.Style)
}
