// eu_ec_eci service for European Citizens' Initiative.
package eu_ec_eci

import (
	"cmp"
	"fmt"
	"slices"
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func Do(t *tool.Tool) {
	eciByYear := fetchAllAcepted(t)
	t.WriteFile("/eu/ec/eci/index.html", lredirect.All)
	t.WriteFile("/eu/ec/eci/schema.html", schemaPage)
	for _, l := range translate.Langs {
		renderIndex(t, eciByYear, l)
	}
	for year := range eciByYear {
		t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/index.html", year), render.Back)
	}
	for _, eciSlice := range eciByYear {
		for _, eci := range eciSlice {
			t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/index.html", eci.Year, eci.Number), lredirect.All) // TODO: the language can be unavailable
			for _, l := range translate.Langs {
				renderOne(t, eci, l)
			}
			if img := eci.Image; img != nil {
				t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/logo%s", eci.Year, eci.Number, img.Raw.Extension), img.Raw.Data)
				if res := img.Resized; res != nil {
					t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/logo%s", eci.Year, eci.Number, res.Extension), res.Data)
				}
			}
		}
	}

	ExploreRefused(t)
}

func fetchAllAcepted(t *tool.Tool) map[int][]*ECIOut {
	checkThreashold(t)

	eciByYear := make(map[int][]*ECIOut)
	for _, info := range fetchAcceptedIndex(t) {
		eci := fetchDetail(t, info)
		if eci == nil {
			continue
		}
		eciByYear[eci.Year] = append(eciByYear[eci.Year], eci)
	}

	for _, byYear := range eciByYear {
		slices.SortFunc(byYear, func(a, b *ECIOut) int {
			return cmp.Compare(b.Number, a.Number)
		})
	}

	return eciByYear
}
