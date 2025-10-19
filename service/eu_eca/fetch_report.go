package eu_eca

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"
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
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/HuguesGuilleus/sniffle/tool/securehtml"
)

type report struct {
	Type    string
	PubDate time.Time

	ImageHash string
	Image     *rimage.Image

	// Avaiable languages for this report.
	Langs [language.Len]bool

	// The textual description of this report.
	// [fetchAnnualReport()] return a report where the description is available for all [translate.Langs].
	Description [language.Len]reportDescription
}

type reportDescription struct {
	// The language of this description
	Lang        language.Language
	Title       string
	Description render.H
	ReportPage  *url.URL
	PDFURL      *url.URL
}

func fetchAnnualReport(t *tool.Tool) (reportByYear map[int][]*report) {
	// Fetch reports.
	mapReports := make(map[string]*report)
	for _, l := range translate.Langs {
		fetchReports(t, mapReports, "Annual report", l)
		fetchReports(t, mapReports, "Special Report", l)
		fetchReports(t, mapReports, "Specific Annual Report", l)
		fetchReports(t, mapReports, "Review", l)
		fetchReports(t, mapReports, "opinions and other outputs", l)
	}

	// Set a description for all languages.
	for _, r := range mapReports {
		baseDesc := cmp.Or(r.Description[language.French], r.Description[language.English])
		for _, l := range translate.Langs {
			if r.Description[l].Title == "" {
				r.Description[l] = baseDesc
			}
		}
	}

	// Sort reports by
	reports := slices.AppendSeq(
		make([]*report, 0, len(mapReports)),
		maps.Values(mapReports),
	)
	slices.SortFunc(reports, func(a, b *report) int {
		return cmp.Or(
			b.PubDate.Truncate(time.Hour*24).Compare(a.PubDate.Truncate(time.Hour*24)),
			cmp.Compare(a.Description[language.French].Title, b.Description[language.French].Title),
		)
	})

	// Group by year.
	reportByYear = make(map[int][]*report)
	for _, r := range reports {
		reportByYear[int(r.PubDate.Year())] = append(reportByYear[int(r.PubDate.Year())], r)
	}

	// [IF DEVMODE] export data in JSON
	if tool.DevMode {
		j, _ := json.MarshalIndent(reportByYear, "", "\t")
		t.WriteFile("/eu/eca/report/src/out.json", j)
	}

	return reportByYear
}

func fetchReports(t *tool.Tool, mapReports map[string]*report, reportType string, l language.Language) {
	request := fetch.Rs(http.MethodPost, "https://www.eca.europa.eu/_vti_bin/ECA.Internet/DocSetService.svc/SearchDocs",
		`{"searchInput":{`+(""+
			`"RowLimit":1000,`+
			`"Filter":"ECADocType=\"`+reportType+`\"",`+
			`"Lang":"`+langName[l]+`",`+
			`"LangCode":"`+l.Upper()+`"}`)+
			`}`,
		"Accept", "application/json",
		"Content-Type", "application/json",
	)

	if tool.DevMode {
		t.WriteFile("/eu/eca/report/src/src."+reportType+"."+l.String()+".json", tool.FetchAll(t, request))
	}

	dto := make([]struct {
		Title                string              `json:"Title"`
		Description          string              `json:"Description"`
		ImageUrl             string              `json:"ImageUrl"`
		PublicationDate      pubDate             `json:"PublicationDate"`
		ReportLandingPageUrl string              `json:"ReportLandingPageUrl"`
		ReportUrl            string              `json:"ReportUrl"`
		Languages            []language.Language `json:"Languages"`
	}, 0)
	if tool.FetchJSON(t, annualReportsType, &dto, request) {
		return
	}

	for _, dto := range dto {
		id := strings.TrimPrefix(strings.Split(dto.ReportLandingPageUrl, "/")[3], "\u200B")
		r := mapReports[id]

		if r == nil {
			r = &report{
				Type:    reportType,
				PubDate: dto.PublicationDate.Time.In(render.DateZone),
			}
			if dto.ImageUrl != "" {
				url := "https://www.eca.europa.eu" + dto.ImageUrl
				h := sha256.Sum256([]byte(url))
				r.Image = rimage.New(t, url)
				r.ImageHash = hex.EncodeToString(h[:5])
			}

			for _, l := range dto.Languages {
				r.Langs[l] = true
			}

			mapReports[id] = r
		}

		r.Langs[l] = true
		r.Description[l] = reportDescription{
			Lang:        l,
			Title:       dto.Title,
			Description: securehtml.Secure(dto.Description),
			ReportPage: &url.URL{
				Scheme: "https",
				Host:   "www.eca.europa.eu",
				Path:   dto.ReportLandingPageUrl,
			},
			PDFURL: securehtml.ParseURL(dto.ReportUrl),
		}
	}
}

var langName = [language.Len]string{
	language.Bulgarian:  "Bulgarian",
	language.Croatian:   "Croatian",
	language.Czech:      "Czech",
	language.Danish:     "Danish",
	language.Dutch:      "Dutch",
	language.English:    "English",
	language.Estonian:   "Estonian",
	language.Finnish:    "Finnish",
	language.French:     "French",
	language.German:     "German",
	language.Greek:      "Greek",
	language.Hungarian:  "Hungarian",
	language.Irish:      "Irish",
	language.Italian:    "Italian",
	language.Latvian:    "Latvian",
	language.Lithuanian: "Lithuanian",
	language.Maltese:    "Maltese",
	language.Polish:     "Polish",
	language.Portuguese: "Portuguese",
	language.Romanian:   "Romanian",
	language.Slovak:     "Slovak",
	language.Slovene:    "Slovene",
	language.Spanish:    "Spanish",
	language.Swedish:    "Swedish",
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
