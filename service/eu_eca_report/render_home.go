package eu_eca_report

import (
	"fmt"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func renderHome(l language.Language, years []int) []byte {
	tr := translate.T[l]
	return render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, "all", tr.EU_ECA.Name, ""),
		render.N("body",
			component.InDevHeader(l),
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerID",
						component.HomeAnchor(l),
						render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
						render.Na("a", "href", l.Path("/eu/eca/")).A("title", tr.EU_ECA.Name).N("eca"),
					),
				),
				render.N("h1", tr.EU_ECA.Name),
				component.HeaderLangs(translate.Langs, l, "/eu/eca/"),
			),
			render.N("main.w",
				render.N("div.edito",
					render.N("div.editoT", tr.GLOBAL.HELP),
					tr.EU_ECA.Desc,
				),
				render.N("div.boxFlex",
					render.Na("a.box", "href", "https://www.eca.europa.eu/"+l.String()+"/").N(tr.GLOBAL.LinkOfficial),
				),
				render.N("div.boxFlex.summary",
					render.Na("a.box.edito", "href", "report/schema.html").N(tr.GLOBAL.SchemaLink),
					render.Na("a.box.data", "href", "report/data.html").N(tr.GLOBAL.ReuseDataLink),
				),
				render.N("div.boxFlex.summary",
					render.Na("a.box", "href", l.Path("report/ar.")).N(tr.EU_ECA.Kind["Annual report"]),
					render.Na("a.box", "href", l.Path("report/rv.")).N(tr.EU_ECA.Kind["Review"]),
					render.Na("a.box", "href", l.Path("report/act.")).N(tr.EU_ECA.Kind["Activity Report"]),
					render.Na("a.box", "href", l.Path("report/j.")).N(tr.EU_ECA.Kind["Journal"]),
					"[!opinions are missing]",
				),
				render.N("div.boxFlex.summary", render.S(years, "", func(year int) render.Node {
					return render.Na("a.box", "href", fmt.Sprintf("report/%d.%s.html", year, l)).N(render.Int(year))
				})),
				render.Na("a.box", "href", l.Path("report/all.")).N(tr.EU_ECA.INDEX_ALL.Name),
			),
			component.Footer(l, 0),
		),
	))
}
