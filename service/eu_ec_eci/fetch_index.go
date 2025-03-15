package eu_ec_eci

import (
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/sch"
)

const acceptedIndexURL = "https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/0/0"

const refusedIndexURL = "https://register.eci.ec.europa.eu/core/api/register/search/REFUSED/EN/0/0"

type indexItem struct {
	id     int
	year   int
	number int
}

// Get all ECI items.
func fetchIndex(t *tool.Tool, ty sch.Type, url string) []indexItem {
	dto := struct {
		Entries []struct {
			ID     int `json:"id"`
			Year   int `json:"year,string"`
			Number int `json:"number,string"`
		} `json:"entries"`
	}{}
	if tool.FetchJSON(t, ty, &dto, fetch.URL(url)) {
		return nil
	}

	items := make([]indexItem, len(dto.Entries))
	for i, dtoEntry := range dto.Entries {
		items[i] = indexItem{
			id:     dtoEntry.ID,
			year:   dtoEntry.Year,
			number: dtoEntry.Number,
		}
	}

	return items
}

// Get all accepted ECI item.
func fetchAcceptedIndex(t *tool.Tool) (items []indexItem) {
	if tool.DevMode {
		t.WriteFile(
			"/eu/ec/eci/src.json",
			tool.FetchAll(t, fetch.URL(acceptedIndexURL)),
		)
	}
	return fetchIndex(t, indexType, acceptedIndexURL)
}

// Get all refused ECI item.
func fetchRefusedIndex(t *tool.Tool) (items []indexItem) {
	if tool.DevMode {
		t.WriteFile(
			"/eu/ec/eci/refused/src.json",
			tool.FetchAll(t, fetch.URL(refusedIndexURL)),
		)
	}
	return fetchIndex(t, nil, refusedIndexURL)
}
