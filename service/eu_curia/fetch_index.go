package eu_curia

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fetch"
	"github.com/HuguesGuilleus/sniffle/tool/sch"
)

func fetchAll(t *tool.Tool) {
	ids := make([]string, 0)
	for year := 1955; year <= time.Now().Year(); year++ {
		t.Warn("fetch.index", "year", year)
		ids = append(ids, fetchIndex(t, requestIndexYear(year))...)
	}

	for i, affID := range ids {
		t.Warn("fetch.procedure", "affID", affID, "%", i*100/len(ids))
		fetchProcedure(t, affID)
	}
}

func fetchMain(t *tool.Tool) {
	ids := make([]string, 0)
	year := time.Now().Year()
	// t.Info("fetch", "year", year-1)
	// ids = append(ids, fetchIndex(t, requestIndexYear(year-1))...)
	t.Info("fetch", "year", year)
	ids = append(ids, fetchIndex(t, requestIndexYear(year))...)
	t.Info("fetch", "state", "pending")
	ids = append(ids, fetchIndex(t, requestIndexPending())...)

	fmt.Println(len(ids))

	for i, affID := range ids {
		_ = i
		// t.Log(context.Background(), slog.LevelInfo+1, "fetch.procedure", "affID", affID, "%", i*100/len(ids))
		fetchProcedure(t, affID)
	}
}

func requestIndexYear(year int) *fetch.Request {
	body, _ := json.Marshal(map[string]any{
		"searchTerm":       "",
		"multiSearchTerms": []any{},
		"sortTermList":     []any{map[string]string{"sortDirection": "DESC", "sortTerm": "SCORE"}},
		"pagination": map[string]any{
			"pageSize":   10_000,
			"pageNumber": 0,
			"from":       0,
			"to":         10_000,
			"origin":     "affair",
		},
		"language":         "FR",
		"tabName":          "affair",
		"isAllTabsRequest": false,
		"ecli":             "",
		"publishedId":      "",
		"usualName":        "",
		"logicDocId":       "",
		"repJurExpand":     true,
		"filtersValue": []map[string]any{
			{
				"field":                   "introductionYear",
				"values":                  []string{strconv.Itoa(year)},
				"valuesWithFullHierarchy": []string{"introductionYear"},
			},
		},
		"isSearchExact": true,
		"searchSources": []string{"metadata"},
	})
	return fetch.R("POST", "https://infocuriaws.curia.europa.eu/elastic-connector/search", body,
		"Content-Type", "application/json; charset=utf-8",
		"Accept", "application/json",
	)
}

func requestIndexPending() *fetch.Request {
	body, _ := json.Marshal(map[string]any{
		"searchTerm":       "",
		"multiSearchTerms": []any{},
		"sortTermList":     []any{map[string]string{"sortDirection": "DESC", "sortTerm": "SCORE"}},
		"pagination": map[string]any{
			"pageSize":   10_000,
			"pageNumber": 0,
			"from":       0,
			"to":         10_000,
			"origin":     "affair",
		},
		"language":         "FR",
		"tabName":          "affair",
		"isAllTabsRequest": false,
		"ecli":             "",
		"publishedId":      "",
		"usualName":        "",
		"logicDocId":       "",
		"repJurExpand":     true,
		"filtersValue": []map[string]any{
			{
				"field":                   "affairState",
				"values":                  []string{"ENC"},
				"valuesWithFullHierarchy": []string{"ENC"},
			},
		},
		"isSearchExact": true,
		"searchSources": []string{"metadata"},
	})
	return fetch.R("POST", "https://infocuriaws.curia.europa.eu/elastic-connector/search", body,
		"Content-Type", "application/json; charset=utf-8",
		"Accept", "application/json",
	)
}

// Fetch Affair index and return Affair ID.
func fetchIndex(t *tool.Tool, request *fetch.Request) []string {
	dto := struct {
		TotalHits  uint `json:"totalHits"`
		SearchHits []struct {
			Content struct {
				AffID string `json:"affId"`
			} `json:"content"`
		} `json:"searchHits"`
	}{}
	if tool.FetchJSON(t, indexTypes, &dto, request) {
		return nil
	}

	ids := make([]string, len(dto.SearchHits))
	for i, hit := range dto.SearchHits {
		ids[i] = hit.Content.AffID
	}
	return ids
}

var indexTypes = sch.Map(
	sch.FieldSR("totalHits", sch.Any()),
	sch.FieldSR("maxScore", sch.Any()),
	sch.FieldSR("maxScoreOthersDocuments", sch.Any()),
	sch.FieldSR("totalHitsRelation", sch.Any()),
	sch.FieldSR("totalHitsOthersDocs", sch.Any()),
	sch.FieldSR("searchHits", sch.Array(sch.Map(
		sch.FieldSR("index", sch.Regexp(`index-affair@\d\d-\d\d-\d\d\d\d`)),
		sch.FieldSR("id", sch.NotEmptyString()),
		sch.FieldSR("score", sch.Or(sch.ConstFloat(0.0), sch.ConstFloat(1.0))),
		sch.FieldSR("sortValues", sch.ArrayItems(
			sch.Or(sch.ConstFloat(0.0), sch.ConstFloat(1.0)),
			sch.AnyInt(),
		)),
		sch.FieldSR("content", sch.MapExtra(
			sch.FieldSR("id", sch.NotEmptyString()),
			sch.FieldSR("affId", sch.Regexp(`[TC]/\d{4}/\d{2}/\w{10}/\d{2}`)),
			sch.FieldSR("publishedId", sch.NotEmptyString()),
			sch.FieldSR("publishedAffId", sch.NotEmptyString()),
			sch.FieldSR("procedureId", sch.NotEmptyString()),
			sch.FieldSR("jurisdiction", sch.EnumString("Cour de justice", "Tribunal")),
			sch.FieldSR("jurisdictionCode", sch.EnumString("T", "C", "F")),
			sch.FieldSR("order", sch.AnyInt()),
			sch.FieldSR("suiteNumber", sch.IntervalInt(1, 3)),
			sch.FieldSR("procType", sch.EnumString("I", "P", "R")),
			sch.FieldSR("procNumber", sch.IntervalInt(1, 4)),
			sch.FieldSR("introductionDate", sch.Time("2006-01-02")),
			sch.FieldSR("affairState", sch.EnumString("Affaire pendante", "Affaire clôturée")),
			sch.FieldSR("affairStateCode", sch.EnumString("ENC", "CLOTPUB")),
			sch.FieldSR("usualNameML", sch.ArrayMin(1, sch.Map(
				sch.Field(languesLowerType, sch.NotEmptyString(), true),
			))),
			sch.FieldSR("joinExist", sch.IntervalInt(0, 1)),
			sch.FieldSR("pilot", sch.IntervalInt(0, 1)),
			sch.FieldSR("joinAffairs", sch.Array(sch.NotEmptyString())),
			sch.FieldSR("pourvoiAffIds", sch.Array(sch.NotEmptyString())),
			sch.FieldSR("originPourvoiIds", sch.Array(sch.NotEmptyString())),
			sch.FieldSR("reExamenAffIds", sch.ArrayMax(1, sch.NotEmptyString())),
			sch.FieldSR("originReExamanIds", sch.ArrayMax(1, sch.NotEmptyString())),
			sch.FieldSR("transfertAffIds", sch.EmptyArray()),
			sch.FieldSR("originTransfertIds", sch.EmptyArray()),
			sch.FieldSR("oppObsFag", sch.ConstInt(0)),
			sch.FieldSR("contentLength", sch.ConstInt(0)),
			sch.FieldSR("docNoPart", sch.ConstInt(0)),
		)),
		sch.FieldSR("highlightFields", sch.Map()),
		sch.FieldSR("innerHits", sch.Map()),
		sch.FieldSR("matchedQueries", sch.EmptyArray()),
	))),
	sch.FieldSR("aggregatesMap", sch.MapExtra()),
)

var languesLowerType = sch.EnumString("bg", "cs", "da", "de", "el", "en", "es", "et", "fi", "fr", "ga", "hr", "hu", "it", "lt", "lv", "mt", "nl", "pl", "pt", "ro", "sk", "sl", "sv")

var languesUpperTypes = sch.EnumString("BG", "CS", "DA", "DE", "EL", "EN", "ES", "ET", "FI", "FR", "GA", "HR", "HU", "IT", "LT", "LV", "MT", "NL", "PL", "PT", "RO", "SK", "SL", "SV")
