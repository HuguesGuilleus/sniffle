package eu_eca_report

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fetch"
)

func fetchReport(t *tool.Tool) {
	types := []string{
		"Annual report", "Specific Annual Report",

		"Review",

		"Activity Report",
		"Journal",
	}

	for _, ty := range types {
		fetchReportType(t, ty, "English", language.English)
		fetchReportType(t, ty, "French", language.French)
	}
}

func fetchReportType(t *tool.Tool, docType, lang string, l language.Language) {
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

	// Basic output
	fmt.Println("fetch", docType, lang, len(dto))
	data := make([]byte, 0)
	for _, report := range dto {
		data = fmt.Appendf(data, "%q /// %q\n- %s\n- %s\n- %s\n\n",
			report.Title, report.Description, report.ReportLandingPageUrl, report.ReportUrl, report.PublicationDate.Time,
		)
	}
	t.WriteFile(
		fmt.Sprintf("/eu/eca/dev.%s.%s.txt", strings.ToLower(strings.ReplaceAll(docType, " ", "_")), l.String()),
		data,
	)
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
