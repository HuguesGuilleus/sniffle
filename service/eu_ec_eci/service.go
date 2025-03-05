package eu_ec_eci

import (
	"cmp"
	"fmt"
	"slices"
	"sniffle/tool"
	"sniffle/tool/render"
)

func Do(t *tool.Tool) {
	eciSlice := fetchAll(t)
	eciByYear := yearGroupingECI(eciSlice)

	t.LangRedirect("/eu/ec/eci/index.html")
	for _, l := range t.Languages {
		renderIndex(t, eciByYear, l)
	}
	renderSchema(t)

	for _, eci := range eciSlice {
		t.LangRedirect(fmt.Sprintf("/eu/ec/eci/%d/%d/index.html", eci.Year, eci.Number))
		for _, l := range t.Languages {
			renderOne(t, eci, l)
		}
		if img := eci.Image; img != nil {
			t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/logo%s", eci.Year, eci.Number, img.Raw.Extension), img.Raw.Data)
			if res := img.Resized; res != nil {
				t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/logo%s", eci.Year, eci.Number, res.Extension), res.Data)
			}
		}
	}

	for y := range eciByYear {
		t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/index.html", y), render.Back)
	}
}

func fetchAll(t *tool.Tool) []*ECIOut {
	checkThreashold(t)

	items := fetchIndex(t)

	eciSlice := make([]*ECIOut, 0, len(items))
	for _, info := range items {
		eci := fetchDetail(t, info)
		if eci == nil {
			continue
		}
		eciSlice = append(eciSlice, eci)
	}

	return eciSlice
}

func yearGroupingECI(eciSlice []*ECIOut) map[int][]*ECIOut {
	slices.SortFunc(eciSlice, func(a, b *ECIOut) int {
		return cmp.Or(
			cmp.Compare(b.Year, a.Year),
			cmp.Compare(b.Number, a.Number),
		)
	})

	eciByYear := make(map[int][]*ECIOut)
	for _, eci := range eciSlice {
		eciByYear[eci.Year] = append(eciByYear[eci.Year], eci)
	}
	return eciByYear
}
