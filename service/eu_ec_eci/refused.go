package eu_ec_eci

import (
	"sniffle/tool"
	"sniffle/tool/fetch"
	"strconv"
)

func ExploreRefused(t *tool.Tool) {
	for _, entry := range fetchRefusedIndex(t) {
		idString := strconv.FormatUint(uint64(entry.id), 10)
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
