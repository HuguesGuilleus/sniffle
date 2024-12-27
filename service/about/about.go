package about

import (
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

// Write about pages.
func Do(t *tool.Tool) {
	basePath := "/about/"
	t.LangRedirect("/about/index.html")

	for _, l := range t.Languages {
		tr := translate.T[l]
		t.WriteFile(l.Path(basePath), render.Merge(render.Na("html", "lang", l.String()).N(
			component.Head(l, t.HostURL+basePath, tr.ABOUT.PageTitle, tr.ABOUT.PageDescription),
			render.N("body.edito",
				component.TopHeader(l),
				render.N("header",
					render.N("div.headerSup", render.N("div.headerId",
						component.HomeAnchor(l), "about",
					)),
					render.N("div.headerTitle", tr.ABOUT.PageTitle),
					component.HeaderLangs(l, ""),
				),
				render.N("div.w", tr.ABOUT.Text),
				component.Footer(l, 0),
			),
		)))
	}
}
