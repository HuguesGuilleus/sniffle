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
	component.RedirectIndex(t, basePath)

	for _, l := range t.Languages {
		tr := translate.AllTranslation[l]
		page := component.Page{
			Language:    l,
			AllLanguage: t.Languages,
			Title:       tr.ABOUT.PageTitle,
			Description: tr.ABOUT.PageDescription,
			BaseURL:     basePath,
			Body: render.N("body.edito",
				component.TopHeader(l),
				component.Header(t.Languages, l, render.N("div.headerId",
					render.No("a",
						render.A("href", "/").A("title", tr.HomeTitle),
						tr.HomeName),
				), render.Z, tr.ABOUT.PageTitle),
				render.N("div.w", tr.ABOUT.Text),
				component.Footer(l),
			),
		}

		component.Html(t, &page)
	}
}
