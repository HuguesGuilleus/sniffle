// home service for /*.html static pages.
package home

import (
	"sniffle/front/component"
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func Do(t *tool.Tool) {
	t.WriteFile("/index.html", lredirect.All)
	for _, l := range translate.Langs {
		tr := translate.T[l]
		t.WriteFile(l.Path("/"), render.Merge(render.Na("html", "lang", l.String()).N(
			component.Head(l, "/", tr.HOME.Name, tr.HOME.PageDescription),
			render.N("body",
				component.InDevHeader(l),
				component.TopHeader(l),
				render.N("header",
					render.N("div.headerTitle", component.HomeAnchor(l), tr.HOME.Name),
					component.HeaderLangs(translate.Langs, l, ""),
				),

				render.N("div.w.home",
					render.N("h1", tr.HOME.EU),
					render.N("ul",
						render.N("li", render.Na("a", "href", l.Path("/eu/ec/eci/")).N("[/eu/ec/eci/] ", tr.EU_EC_ECI.Name)),
						render.IfS(tool.DevMode, render.N("li", render.Na("a", "href", l.Path("/eu/eca/report/")).N("[/eu/eca/report] ", "!ECA report"))),
					),

					render.N("h1", tr.HOME.About),
					render.N("ul",
						render.N("li", render.Na("a", "href", l.Path("/about/")).N("[/about/] ", tr.ABOUT.PageTitle)),
						render.N("li", render.Na("a", "href", "/release/").N("[/release/] ", tr.HOME.Release)),
					),
				),

				component.Footer(l, 0),
			),
		)))
	}
}
