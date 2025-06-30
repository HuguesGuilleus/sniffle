// about service for /about/*.html static pages.
package about

import (
	"sniffle/front/component"
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

// Write about pages.
func Do(t *tool.Tool) {
	t.WriteFile("/about/index.html", lredirect.All)

	basePath := "/about/"
	for _, l := range translate.Langs {
		tr := translate.T[l]
		t.WriteFile(l.Path(basePath), render.Merge(render.Na("html", "lang", l.String()).N(
			component.Head(l, basePath, tr.ABOUT.PageTitle, tr.ABOUT.PageDescription),
			render.N("body.edito",
				component.TopHeader(l),
				render.N("header",
					render.N("div.headerSup", render.N("div.headerID",
						component.HomeAnchor(l), "about",
					)),
					render.N("div.headerTitle", tr.ABOUT.PageTitle),
					component.HeaderLangs(translate.Langs, l, ""),
				),
				render.N("div.w",
					render.S(tr.ABOUT.Intro, "", func(p render.H) render.Node {
						return render.N("p", p)
					}),
					render.N("hr"),
					render.N("p.noindent", tr.ABOUT.Mail, "ghugues[at]netc[dot]fr"),
					render.N("p.noindent", tr.ABOUT.Host, "OVH SAS\u202F; 2 rue Kellermann, 59100 Roubaix. France."),
				),
				component.Footer(l, 0),
			),
		)))
	}
}
