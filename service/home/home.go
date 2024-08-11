package home

import (
	"context"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func Do(_ context.Context, t *tool.Tool) {
	t.LangRedirect("/index.html")
	for _, l := range t.Languages {
		tr := translate.AllTranslation[l]
		component.Html(t, &component.Page{
			Language:    l,
			AllLanguage: t.Languages,
			Title:       "Home",
			Description: tr.HOME.PageDescription,
			BaseURL:     "/",
			Body: render.N("body ",
				component.TopHeader(l),
				component.Header(t.Languages, l, render.Z, render.Z, tr.HOME.Name),
				render.N("ul.w",
					render.N("li", render.No("a", render.A("href", "/about/"+l.String()+".html"), tr.ABOUT.PageTitle)),
					render.N("li", render.No("a", render.A("href", "/eu/ec/eci/"+l.String()+".html"), tr.EU_EC_ECI.Name)),
				),
				component.Footer(l),
			),
		})
	}
}
