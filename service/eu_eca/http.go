package eu_eca

import (
	"net/http"
	"net/url"
	"slices"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"sniffle/tool/securehtml"
	"time"
)

// (ECADocType= "Annual report" OR ECADocType= "Review" OR ECADocType= "Special Report" OR ECADocType= "Specific Annual Report" OR ECADocType= "opinions and other outputs")

func Do(t *tool.Tool) {

	// if tool.DevMode {
	// 	t.WriteFile("/eu/eca/src-token.json", tool.FetchAll(t,
	// 		fetch.R(http.MethodPost, "https://www.eca.europa.eu/_api/contextinfo", nil,
	// 			"Accept", "application/json; odata=verbose",
	// 			"Content-Type", "application/json; odata=verbose",
	// 		),
	// 	))
	// }
	// v := struct {
	// 	D struct {
	// 		GetContextWebInformation struct {
	// 			FormDigestValue string
	// 		}
	// 	}
	// }{}
	// tool.FetchJSON(t, nil, &v, fetch.R(http.MethodPost, "https://www.eca.europa.eu/_api/contextinfo", nil,
	// 	"Accept", "application/json; odata=verbose",
	// 	"Content-Type", "application/json; odata=verbose",
	// ))
	// token := v.D.GetContextWebInformation.FormDigestValue

	t.LangRedirect("/eu/eca/annual-report/index.html")
	reports := fetchReports(t, language.French)
	for _, l := range t.Languages {
		renderReportIndex(t, l, reports)
	}
}

type Report struct {
	Title       string
	Description render.H
	Publication time.Time
	ReportPage  *url.URL
	ReportURL   *url.URL
	Languages   []language.Language
}

func fetchReports(t *tool.Tool, lang language.Language) []Report {
	langName := "French"
	request := fetch.Rs(http.MethodPost, "https://www.eca.europa.eu/_vti_bin/ECA.Internet/DocSetService.svc/SearchDocs", `{"searchInput":{"RowLimit":1000,"Filter":"ECADocType=\"Annual report\"","Lang":"`+langName+`","LangCode":"FR"}}`,
		"Accept", "application/json",
		"Content-Type", "application/json",
	)

	if tool.DevMode {
		t.WriteFile("/eu/eca/src/doc."+lang.String()+".json", tool.FetchAll(t, request))
	}

	dto := make([]struct {
		Title                string
		Description          string
		ImageUrl             string
		PublicationDate      pubDate
		ReportLandingPageUrl string
		ReportUrl            string
		Languages            []language.Language
	}, 0)
	if tool.FetchJSON(t, annualReportsType, &dto, request) {
		return nil
	}

	reports := make([]Report, len(dto))
	for i, dto := range dto {
		slices.Sort(dto.Languages)
		reports[i] = Report{
			Title:       dto.Title,
			Description: securehtml.Secure(dto.Description),
			Publication: dto.PublicationDate.Time,
			ReportPage:  securehtml.ParseURL("https://www.eca.europa.eu" + dto.ReportLandingPageUrl),
			ReportURL:   securehtml.ParseURL(dto.ReportUrl),
			Languages:   dto.Languages,
			// Image
		}
	}

	return reports
}

type pubDate struct {
	Time time.Time
}

func (p *pubDate) UnmarshalText(text []byte) error {
	t, err := time.Parse("1/2/2006 15:04:05 PM", string(text))
	if err != nil {
		return err
	}
	p.Time = t
	return nil
}
