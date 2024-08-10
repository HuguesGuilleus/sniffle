package eu_ec_eci

import (
	"context"
	"fmt"
	"sniffle/front/component"
	"sniffle/tool"
)

func Do(ctx context.Context, t *tool.Tool) {
	eciSlice := fetchAll(ctx, t)

	component.RedirectIndex(t, "/eu/ec/eci/")

	for _, eci := range eciSlice {
		component.RedirectIndex(t, fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number))
		for _, l := range t.Languages {
			renderOne(t, eci, l)
		}
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
