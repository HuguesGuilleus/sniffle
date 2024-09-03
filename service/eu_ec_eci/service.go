package eu_ec_eci

import (
	"context"
	"fmt"
	"sniffle/tool"
	"sniffle/tool/render"
)

func Do(ctx context.Context, t *tool.Tool) {
	eciSlice := fetchAll(ctx, t)

	t.LangRedirect("/eu/ec/eci/index.html")
	for _, l := range t.Languages {
		renderIndex(t, eciSlice, l)
	}

	years := make(map[int]bool)
	for _, eci := range eciSlice {
		years[eci.Year] = true
		t.LangRedirect(fmt.Sprintf("/eu/ec/eci/%d/%d/index.html", eci.Year, eci.Number))
		for _, l := range t.Languages {
			renderOne(t, eci, l)
		}
		if eci.ImageName != "" {
			t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/%s", eci.Year, eci.Number, eci.ImageName), eci.ImageData)
		}
	}

	for y := range years {
		t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/index.html", y), render.Back)
	}
}

func fetchAll(ctx context.Context, t *tool.Tool) []*ECIOut {
	items := fetchIndex(ctx, t)

	eciSlice := make([]*ECIOut, 0, len(items))
	for _, info := range items {
		eci := fetchDetail(ctx, t, info)
		if eci == nil {
			continue
		}
		eciSlice = append(eciSlice, eci)
	}

	return eciSlice
}
