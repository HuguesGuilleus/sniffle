package eu_eca

import (
	"slices"
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
	"strconv"
	"strings"
)

func renderReportIndex(t *tool.Tool, l language.Language, reportByYear map[int][]*report) {
	tr := translate.T[l]
	basePath := "/eu/eca/report/"
	t.WriteFile(l.Path(basePath), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, basePath, "$ECA report", "$ECA report index"),
		render.N("body",
			component.InDevHeader(l),
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					headerSupECA(l),
					render.N("div.headerID", "report"),
				),
				render.N("div.headerTitle", "$ECA report index"),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.wt", component.Toc(l), render.N("div.wc",
				render.N("div.summary",
					render.Na("a.box", "href", "schema.html").N(tr.GLOBAL.SchemaLink),
				),
				component.SearchBlock(l),
				render.MapReverse(reportByYear, func(year int, reports []*report) render.Node {
					return render.N("div.sg",
						render.Na("h1", "id", strconv.Itoa(year)).N(render.Int(year)),
						render.S(reports, "", func(r *report) render.Node {
							desc := r.Description[l]
							main := render.N("",
								render.N("div.itemTitle.st", desc.Title),
								render.IfS(desc.Description != "", render.N("div.itemDesc", desc.Description)),
							)

							return render.N("div.si.bigItem",
								render.N("div",
									render.N("span.tag", r.PubDate),
									render.N("span.tag", "$$$", r.Type),
								),
								render.IfS(r.Langs[l], main),
								render.If(!r.Langs[l], func() render.Node {
									return render.Na("div.tr", "lang", l.String()).N(main)
								}),
								render.If(desc.ReportPage != nil || desc.PDFURL != nil, func() render.Node {
									return render.N("div.boxFlex",
										render.If(desc.ReportPage != nil, func() render.Node {
											return render.Na("a.box", "href", desc.ReportPage.String()).N("$ page ~>")
										}),
										render.If(desc.PDFURL != nil, func() render.Node {
											return render.Na("a.doc", "href", desc.PDFURL.String()).N("$ pdf ~>")
										}),
									)
								}),
								render.N("div", "$avaiable langs: ", renderAvailableLangs(l, r.Langs)),
								r.Image.Render(r.ImageHash, ""),
							)
						}),
					)
				})),
			),
			component.Footer(l, component.JsSearch|component.JsToc),
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

func renderAvailableLangs(pageLang language.Language, langs [language.Len]bool) string {
	tr := translate.T[pageLang]
	available := make([]string, 0, language.Len)
	for l, ok := range langs {
		if ok {
			available = append(available, string(tr.Langage[l]))
		}
	}
	slices.Sort(available)
	return strings.Join(available, ", ")
}
