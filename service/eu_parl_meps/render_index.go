package eu_parl_meps

import (
	"fmt"
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func renderIndex(t *tool.Tool, l language.Language, list []meps) {
	tr := translate.T[l]
	baseURL := "/eu/parl/meps/"

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, baseURL, "!list eurodéputés ", "!liste basic des eurodéputées"),
		render.N("body",
			component.InDevHeader(l),
			component.TopHeader(l),

			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerID",
						component.HomeAnchor(l),
						render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
						render.Na("a", "href", l.Path("/eu/parl/")).A("title", "!parelement européen").N("parl"), " / ",
						render.Na("a", "href", l.Path("/eu/parl/meps/")).A("title", "!Member of european parlement").N("meps"),
					),
				),
				render.N("div.headerTitle", "!Liste eurodéputées"),
				component.HeaderLangs(translate.Langs, l, ""),
			),

			render.N("main.w",
				render.N("ul",
					render.S(list, "", func(m meps) render.Node {
						return render.N("li", render.Na("a", "href", fmt.Sprintf("%d/%s.html", m.Identifier, l)).N(
							m.Label, " ",
							render.S(m.Term, " ", func(term uint) render.Node {
								return render.N("", "!term:", term)
							}),
							" #", render.Int(m.Identifier),
						))
					}),
				),
				render.N(""),
			),

			component.Footer(l, 0),
		),
	)))
}
