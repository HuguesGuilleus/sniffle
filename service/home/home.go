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
		component.HTML(t, &component.Page{
			Language:    l,
			Title:       "Home",
			Description: tr.HOME.PageDescription,
			BaseURL:     "/",
		}, render.N("body",
			component.TopHeader(l),
			component.Header(t.Languages, l, render.Z, render.Z, tr.HOME.Name),
			render.No("ul.w",
				render.A("style", "font-size:xx-large"),
				render.N("li", render.No("a", render.A("href", "/about/"+l.String()+".html"), tr.ABOUT.PageTitle)),
				render.N("li", render.No("a", render.A("href", "/eu/ec/eci/"+l.String()+".html"), tr.EU_EC_ECI.Name)),
			),
			component.Footer(l),
		))
	}
}
