package component

import (
	"encoding/json"
	"fmt"

	"github.com/HuguesGuilleus/sniffle/tool/render"
)

// Print the v in JSON, and render in a .working component.
func Json(v any) render.Node {
	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Println("!!!!! [component.Json]", err.Error())
		return render.N("pre.working", string(err.Error()))
	}
	return render.N("div.working", render.N("pre", render.H(j)))
}
