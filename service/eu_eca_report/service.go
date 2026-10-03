package eu_eca_report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/HuguesGuilleus/sniffle/front/lredirect"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func Do(t *tool.Tool) {
	allReports := fetchAllReport(t)
	reports := allReports.All()

	t.WriteFile("/eu/eca/report/index.txt", []byte(strings.Join(allReports.Ids(), "\n")))
	t.WriteFile("/eu/eca/report/index.html", render.Back)
	t.WriteFile("/eu/eca/index.html", lredirect.All)
	t.WriteFile("/eu/eca/report/schema.html", schemaPage)

	for _, report := range reports {
		j, _ := json.Marshal(report)
		t.WriteFile("/eu/eca/report/"+report.ID+".json", j)
		report.Image.Save(t, "/eu/eca/report/"+report.ID)
	}
	for _, l := range translate.Langs {
		t.WriteFile(l.Path("/eu/eca/report/all."), renderAll(l, reports))
		t.WriteFile(l.Path("/eu/eca/"), renderHome(l, getYearsRange(reports)))
		t.WriteFile(l.Path(pathByKind("Annual report")), renderByKind(l, "Annual report", allReports.AnnualReport))
		t.WriteFile(l.Path(pathByKind("Review")), renderByKind(l, "Review", allReports.Review))
		t.WriteFile(l.Path(pathByKind("Activity Report")), renderByKind(l, "Activity Report", allReports.ActivityJournal))
		t.WriteFile(l.Path(pathByKind("Journal")), renderByKind(l, "Journal", allReports.Journal))
		for year, reports := range allReports.ByYear() {
			t.WriteFile(
				fmt.Sprintf("/eu/eca/report/%d.%s.html", year, l),
				renderByYear(l, year, reports),
			)
		}
	}
}

func pathByKind(kind string) string {
	switch kind {
	case "Annual report":
		return "/eu/eca/report/ar."
	case "Specific Annual Report":
		return "/eu/eca/report/sar."
	case "Review":
		return "/eu/eca/report/rv."
	case "Activity Report":
		return "/eu/eca/report/act."
	case "Journal":
		return "/eu/eca/report/j."
	}
	return "?"
}
