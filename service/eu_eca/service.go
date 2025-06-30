// eu_eca is service for European Court of Auditors (composition and reports).
package eu_eca

import (
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
)

func Do(t *tool.Tool) {
	reportByYear := fetchAnnualReport(t)

	t.WriteFile("/eu/eca/report/schema.html", schemaPage)
	t.WriteFile("/eu/eca/report/index.html", lredirect.All)

	for _, l := range translate.Langs {
		renderReportIndex(t, l, reportByYear)
	}

	for _, reports := range reportByYear {
		for _, r := range reports {
			r.Image.Save(t, "/eu/eca/report/"+r.ImageHash)
		}
	}
}
