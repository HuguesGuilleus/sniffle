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
		}
		page.Body = render.N("body.edito",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerId",
						render.No("a",
							render.A("href", "/").A("title", tr.HomeTitle),
							tr.HomeName),
					),
				),
				render.N("div.headerTitle", tr.ABOUT.PageTitle),
			),
			render.N("div.w", tr.ABOUT.Text),
			component.Footer(l),
		)

		component.Html(t, &page)
	}
}
