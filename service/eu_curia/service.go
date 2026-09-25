package eu_curia

import "github.com/HuguesGuilleus/sniffle/tool"

func Do(t *tool.Tool) {
	t.WriteFile("/eu/curia/schema.html", schemaHTML)

	fetchMain(t)
}
