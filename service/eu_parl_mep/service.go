package eu_parl_mep

import (
	"fmt"
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
)

const termLen = 11

func Do(t *tool.Tool) {
	terms := make([][]mep, termLen)
	for term := range terms {
		terms[term] = fetchMep(t, term)
	}

	list := make([]mep, 0)
	for _, terms := range terms {
		list = append(list, terms...)
	}

	t.WriteFile("/eu/parl/mep/schema.html", schemaPage)
	t.WriteFile("/eu/parl/mep/index.html", lredirect.All)
	for _, l := range translate.Langs {
		renderIndex(t, l)
	}

	for termIndex, termMeps := range terms {
		t.WriteFile(fmt.Sprintf("/eu/parl/mep/term-%d/index.html", termIndex), lredirect.All)
		for _, l := range translate.Langs {
			renderTerm(t, l, termIndex, termMeps)
		}
	}

	for _, m := range list {
		t.WriteFile(fmt.Sprintf("/eu/parl/mep/%d/index.html", m.Identifier), lredirect.All)
		for _, l := range translate.Langs {
			renderMeps(t, l, m)
		}
	}
}
