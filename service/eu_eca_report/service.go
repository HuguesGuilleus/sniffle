package eu_eca_report

import (
	"encoding/json"
	"strings"

	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func Do(t *tool.Tool) {
	reports := fetchAllReport(t)

	t.WriteFile("/eu/eca/report/index.txt", []byte(strings.Join(reports.Ids(), "\n")))

	for report := range reports.All() {
		j, _ := json.Marshal(report)
		t.WriteFile("/eu/eca/report/"+report.ID+"/data.json", j)
	}

	t.WriteFile("/eu/eca/report/", render.Back)
}
