// eu_eca is service for European Court of Auditors (composition and reports).
package eu_eca

import (
	"fmt"

	"github.com/HuguesGuilleus/sniffle/front/lredirect"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
)

func Do(t *tool.Tool) {
	reportByYear := fetchAnnualReport(t)

	t.WriteFile("/eu/eca/report/schema.html", schemaPage)
	t.WriteFile("/eu/eca/report/index.html", lredirect.All)

	for year, reports := range reportByYear {
		t.WriteFile(fmt.Sprintf("/eu/eca/%d/index.html", year), lredirect.All)
		for _, l := range translate.Langs {
			renderIndexByYear(t, l, year, reports)
		}
	}

	for _, reports := range reportByYear {
		for _, r := range reports {
			r.Image.Save(t, "/eu/eca/report/"+r.ImageHash)
		}
	}
}
