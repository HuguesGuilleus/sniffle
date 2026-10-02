package eu_eca_report

import (
	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func renderAll(l language.Language, reports []report) []byte {
	tr := translate.T[l]

	return render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, "all", tr.EU_ECA.INDEX_ALL.Name, tr.EU_ECA.INDEX_ALL.Desc),
		render.N("body",
			component.InDevHeader(l),
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerID",
						component.HomeAnchor(l),
						render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
						render.Na("a", "href", l.Path("/eu/eca/")).A("title", tr.EU_ECA.Name).N("eca"), " / ",
						render.Na("a", "href", l.Path("/eu/eca/report/")).N("report")),
					render.N("div.headerID", "all"),
					render.N("h1", tr.EU_ECA.INDEX_ALL.Name),
				),
				component.HeaderLangs(translate.Langs, l, "all."),
			),
			render.N("main.ww",
				component.SearchBlock(l),
				render.N("ul", render.S(reports, "", func(r report) render.Node {
					rl := r.L[l]
					langingPage := render.If(rl.ReportLandingPageUrl != nil, func() render.Node {
						return render.Na("a", "href", rl.ReportLandingPageUrl.String()).N(tr.EU_ECA.ReportPage, " ~>")
					})
					pdfURL := render.If(rl.ReportLandingPageUrl != nil, func() render.Node {
						return render.Na("a", "href", r.L[l].ReportLandingPageUrl.String()).N(tr.EU_ECA.ReportPDF, " ~>")
					})
					title := render.N("", r.L[l].Title, " ", langingPage, " ", pdfURL)
					if r.L[l].L != l {
						title = render.Na("span.tri", "lang", r.L[l].L.String()).N(r.L[l].Title, " ", langingPage, " ", pdfURL)
					}
					return render.N("li.si",
						tr.DateShort(r.PublicationDate),
						" | ", tr.EU_ECA.Kind[r.Kind],
						" [", r.ID, "] ",

						title,
					)
				})),
			),
			component.Footer(l, component.JsSearch),
		),
	))
}
