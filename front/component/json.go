package component

import (
	"encoding/json"
	"fmt"
	"sniffle/tool/render"
)

func Json(v any) render.Node {
	j, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Println("!!!!! [component.Json]", err.Error())
		return render.N("pre.working", string(err.Error()))
	}
	return render.N("div.working", render.N("pre", render.H(j)))
}
