package eu_eca_report

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/common/rimage"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fetch"
	"github.com/HuguesGuilleus/sniffle/tool/securehtml"
)

func fetchAllReport(t *tool.Tool) allReport {
	return allReport{
		AnnualReport:         fetchReportsByType(t, "Annual report"),
		SpecificAnnualReport: fetchReportsByType(t, "Specific Annual Report"),
		Review:               fetchReportsByType(t, "Review"),
		ActivityJournal:      fetchReportsByType(t, "Activity Report"),
		Journal:              fetchReportsByType(t, "Journal"),
	}
}

func fetchReportsByType(t *tool.Tool, ty string) (all []report) {
	langs := [language.Len]string{
		language.English: "English",
		language.French:  "French",
	}
	for _, l := range translate.Langs {
		all = mergeReports(all, fetchReportType(t, ty, langs[l], l))
	}
	fillDefault(all)
	return
}

func fetchReportType(t *tool.Tool, docType, lang string, l language.Language) (reports []report) {
	body, _ := json.Marshal(struct {
		S map[string]any `json:"searchInput"`
	}{
		S: map[string]any{
			"RowLimit": 10_000,

			"Filter": `ECADocType="` + docType + `"`,

			"Lang":     lang,
			"LangCode": l.Upper(),
		},
	})

	dto := []struct {
		Title       string `json:"Title"`
		Description string `json:"Description"`

		ImageUrl        string  `json:"ImageUrl"`
		PublicationDate timeDTO `json:"PublicationDate"`

		ReportLandingPageUrl string   `json:"ReportLandingPageUrl"`
		ReportUrl            string   `json:"ReportUrl"`
		Langues              []string `json:"Languages"`
	}{}

	request := fetch.R(http.MethodPost, "https://www.eca.europa.eu/_vti_bin/ECA.Internet/DocSetService.svc/SearchDocs", body,
		"Accept", "application/json",
		"Content-Type", "application/json",
		"Referer", "https://www.eca.europa.eu/"+l.String()+"/multiple-reports",
	)
	if tool.DevMode {
		t.WriteFile(
			fmt.Sprintf("/eu/eca/dev.%s.%s.json", strings.ToLower(strings.ReplaceAll(docType, " ", "_")), l.String()),
			tool.FetchAll(t, request),
		)
	}
	if tool.FetchJSON(t, indexTypes, &dto, request) {
		t.Logger.Warn("json.fail", "type", docType, "lang", lang)
		return
	}

	reports = make([]report, len(dto))
	for i, r := range dto {
		reports[i] = report{
			ID:              r.ReportLandingPageUrl[len("/../publications/"):],
			Kind:            docType,
			PublicationDate: r.PublicationDate.Time.UTC(),
		}
		if r.ImageUrl != "" {
			url := "https://www.eca.europa.eu" + r.ImageUrl
			reports[i].ImageURL = url
			reports[i].Image = rimage.New(t, url)
		}
		reports[i].L[l] = &reportL{
			L:                    l,
			Title:                r.Title,
			Description:          securehtml.Secure(r.Description),
			ReportLandingPageUrl: &url.URL{Scheme: "https", Host: "www.eca.europa.eu", Path: r.ReportLandingPageUrl},
			ReportUrl:            securehtml.ParseURL(r.ReportUrl),
		}
	}

	slices.SortFunc(reports, func(a, b report) int {
		return b.PublicationDate.Compare(a.PublicationDate)
	})

	return reports
}

func mergeReports(reportsListA, reportsListB []report) (merge []report) {
	if len(reportsListA) == 0 {
		return reportsListB
	} else if len(reportsListB) == 0 {
		return reportsListA
	}
	merge = make([]report, 0, len(reportsListA)+len(reportsListB))
	a, b := 0, 0
	for a < len(reportsListA) && b < len(reportsListB) {
		if reportsListA[a].ID == reportsListB[b].ID {
			out := reportsListA[a]
			for l, reportL := range reportsListB[b].L {
				if reportL != nil {
					out.L[l] = reportL
				}
			}
			merge = append(merge, out)
			a++
			b++
		} else if reportsListA[a].ID < reportsListB[b].ID {
			merge = append(merge, reportsListA[a])
			a++
		} else {
			merge = append(merge, reportsListB[b])
			b++
		}
	}
	merge = append(merge, reportsListA[a:]...)
	merge = append(merge, reportsListB[b:]...)
	return
}

func fillDefault(reports []report) {
	for i := range reports {
		for l := range reports[i].L {
			reports[i].L[l] = cmp.Or(
				reports[i].L[l],
				reports[i].L[language.English],
				reports[i].L[language.French],
			)
		}
	}
}

type timeDTO struct {
	Time time.Time
}

func (dto *timeDTO) UnmarshalText(data []byte) error {
	t, err := time.Parse("1/2/2006 15:04:05 PM", string(data))
	if err != nil {
		return err
	}
	dto.Time = t
	return nil
}
