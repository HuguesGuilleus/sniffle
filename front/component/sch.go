package component

import "github.com/HuguesGuilleus/sniffle/tool/render"

func SchComment(args ...any) render.Node {
	return render.N("", render.N("span.sch-comment", "// ", render.N("", args...)), "\n")
}
