package eu_ec_eci

import (
	"sniffle/tool"
	"sniffle/tool/fetch"
	"strconv"
)

const refusedIndexURL = "https://register.eci.ec.europa.eu/core/api/register/search/REFUSED/FR/0/0"

func ExploreRefused(t *tool.Tool) {
	if tool.DevMode {
		t.WriteFile("/eu/ec/eci/refused/src.json", tool.FetchAll(t, fetch.URL(refusedIndexURL)))
	}
	dtoIndex := struct{ Entries []struct{ ID int } }{}
	if tool.FetchJSON(t, refusedIndexType, &dtoIndex, fetch.URL(refusedIndexURL)) {
		return
	}

	for _, entry := range dtoIndex.Entries {
		idString := strconv.Itoa(entry.ID)
		request := fetch.URL("https://register.eci.ec.europa.eu/core/api/register/details/" + idString)
		if tool.DevMode {
			t.WriteFile("/eu/ec/eci/refused/"+idString+"/src.json", tool.FetchAll(t, request))
		}
		dtoOne := struct{}{}
		if tool.FetchJSON(t, refusedOneType, &dtoOne, request) {
			return
		}
	}
}
