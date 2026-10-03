package eu_eca_report

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func renderByYear(l language.Language, year int, reports []*report) []byte {
	tr := translate.T[l]
	title := fmt.Sprintf(tr.EU_ECA.INDEX_BY_YEAR.Name, year)
	return render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, "all", title, tr.EU_ECA.INDEX_BY_YEAR.Desc),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerID",
						component.HomeAnchor(l),
						render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
						render.Na("a", "href", l.Path("/eu/eca/")).A("title", tr.EU_ECA.Name).N("eca"),
					),
					render.N("div.headerID", "report / ", render.Int(year)),
				),
				render.N("h1", title),
				component.HeaderLangs(translate.Langs, l, strconv.Itoa(year)+"."),
			),
			render.N("main.w",
				component.SearchBlock(l),
				render.S(reports, "", func(r *report) render.Node { return renderBigItem(l, r) }),
			),
			component.Footer(l, component.JsSearch),
		),
	))
}

func renderByKind(l language.Language, kind string, reports []report) []byte {
	tr := translate.T[l]
	return render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, "all", tr.EU_ECA.Kind[kind], tr.EU_ECA.INDEX_BY_KIND.Desc+tr.EU_ECA.Kind[kind]),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerID",
						component.HomeAnchor(l),
						render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
						render.Na("a", "href", l.Path("/eu/eca/")).A("title", tr.EU_ECA.Name).N("eca"),
					),
					render.N("div.headerID", "report / ", strings.ToLower(kind)),
				),
				render.N("h1", tr.EU_ECA.Kind[kind]),
				component.HeaderLangs(translate.Langs, l, pathByKind(kind)),
			),
			render.N("main.w",
				render.N("div.edito",
					render.N("div.editoT", tr.GLOBAL.HELP),
					tr.EU_ECA.INDEX_BY_KIND.Edito[kind],
				),
				component.SearchBlock(l),
				render.S(reports, "", func(r report) render.Node { return renderBigItem(l, &r) }),
			),
			component.Footer(l, component.JsSearch),
		),
	))
}

func renderBigItem(l language.Language, r *report) render.Node {
	tr := translate.T[l]
	rl := r.L[l]
	tag := "div.bigItem.si"
	if rl.L != l {
		tag += ".tr"
	}
	return render.N(tag,
		r.Image.Render("/eu/eca/report/"+r.ID, tr.GLOBAL.LogoTitle),
		render.N("div.itemTitle.st", rl.Title),
		render.N("div", tr.EU_ECA.PublicationDate, tr.DateLong(r.PublicationDate)),
		render.If(rl.Description != "", func() render.Node {
			return render.N("div.itemDesc.st", rl.Description)
		}),
		render.N("div.boxFlex",
			render.N("span.box", r.ID),
			render.If(rl.ReportLandingPageUrl != nil, func() render.Node {
				return render.Na("a.box", "href", rl.ReportLandingPageUrl.String()).N(tr.EU_ECA.ReportPage, " ~>")
			}),
			render.If(rl.ReportLandingPageUrl != nil, func() render.Node {
				return render.Na("a.box", "href", r.L[l].ReportLandingPageUrl.String()).N(tr.EU_ECA.ReportPDF, " ~>")
			}),
		),
	)
}
