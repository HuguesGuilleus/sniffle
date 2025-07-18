package eu_parl_meps

import (
	"fmt"
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
)

func Do(t *tool.Tool) {
	list := fetchMeps(t)

	t.WriteFile("/eu/parl/meps/index.html", lredirect.All)
	for _, l := range translate.Langs {
		renderIndex(t, l, list)
	}

	for _, m := range list {
		t.WriteFile(fmt.Sprintf("/eu/parl/meps/%d/index.html", m.Identifier), lredirect.All)
		for _, l := range translate.Langs {
			renderMeps(t, l, m)
		}
	}
}
