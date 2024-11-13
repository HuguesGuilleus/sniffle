package home

import (
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func Do(t *tool.Tool) {
	t.LangRedirect("/index.html")
	for _, l := range t.Languages {
		tr := translate.T[l]

		t.WriteFile(l.Path("/"), render.Merge(render.Na("html", "lang", l.String()).N(
			component.Head(l, "/", tr.HOME.Name, tr.HOME.PageDescription),
			render.N("body",
				component.TopHeader(l),
				render.N("header",
					render.N("div.headerTitle", component.HomeAnchor(l), tr.HOME.Name),
					component.HeaderLangs(l, ""),
				),
				render.Na("ul.w", "style", "font-size:xx-large").N(
					render.N("li", render.Na("a", "href", l.Path("/about/")).N(tr.ABOUT.PageTitle)),
					render.N("li", render.Na("a", "href", l.Path("/eu/ec/eci/")).N(tr.EU_EC_ECI.Name))),
				component.Footer(l, 0),
			),
		)))
	}
}
