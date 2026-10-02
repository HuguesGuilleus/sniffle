package eu_eca_report

import (
	"encoding/json"
	"strings"

	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func Do(t *tool.Tool) {
	allReports := fetchAllReport(t)

	t.WriteFile("/eu/eca/report/index.txt", []byte(strings.Join(allReports.Ids(), "\n")))
	t.WriteFile("/eu/eca/report/index.html", render.Back)

	reports := allReports.All()
	for _, report := range reports {
		j, _ := json.Marshal(report)
		t.WriteFile("/eu/eca/report/"+report.ID+".json", j)
	}
	for _, l := range translate.Langs {
		t.WriteFile(l.Path("/eu/eca/report/all."), renderAll(l, reports))
	}
}
