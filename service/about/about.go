package about

import (
	"context"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

// Write about pages.
func Do(_ context.Context, t *tool.Tool) {
	basePath := "/about/"
	t.LangRedirect("/about/index.html")

	for _, l := range t.Languages {
		tr := translate.AllTranslation[l]
		component.HTML(t, &component.Page{
			Language:    l,
			Title:       tr.ABOUT.PageTitle,
			Description: tr.ABOUT.PageDescription,
			BaseURL:     basePath,
		}, render.N("body.edito",
			component.TopHeader(l),
			component.Header(t.Languages, l,
				render.N("div.headerId", component.HomeAnchor(l)),
				render.Z, tr.ABOUT.PageTitle),
			render.N("div.w", tr.ABOUT.Text),
			component.Footer(l, 0),
		))
	}
}
