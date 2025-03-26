package eu_eca

import (
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func renderReportIndex(t *tool.Tool, l language.Language, reports []Report) {
	basePath := "/eu/eca/annual-report/"
	hostURL := t.HostURL + basePath
	t.WriteFile(l.Path(basePath), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, hostURL, "$ECA report", "$ECA report index"),
		render.N("body",
			component.TopHeader(l),
			component.InDevHeader(l),
			render.N("header",
				render.N("div.headerSup",
					headerSupECA(l),
					render.N("div.headerID", "annual-report"),
				),
				render.N("div.headerTitle", "$ECA report index"),
				// traductions...
			),
			render.N("main.w",
				render.N("ul", render.S(reports, "", func(r Report) render.Node {
					return render.N("li.doc",
						render.N("div.docT", r.Title),
						render.N("div", "$Publication: ", r.Publication),
						render.N("div", "$Langues: ", render.S(r.Languages, ", ", func(l language.Language) render.Node {
							return render.N("", l.String())
						})),
						render.N("div",
							render.Na("a.doc", "href", r.ReportPage.String()).N("$Page"),
							" ",
							render.Na("a.doc", "href", r.ReportURL.String()).N("$Rapport"),
						),
						render.N("p", r.Description),
					)
				})),
			),
			component.Footer(l, component.JsSearch),
		),
	)))
}

func headerSupECA(l language.Language) render.Node {
	tr := translate.T[l]
	return render.N("div.headerID",
		component.HomeAnchor(l),
		render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
		render.Na("a", "href", l.Path("/eu/eca/")).A("title", "$Europea court of Auditor").N("eca"),
	)
}
