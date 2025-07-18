package eu_parl_meps

import (
	"fmt"
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func renderMeps(t *tool.Tool, l language.Language, m meps) {
	tr := translate.T[l]
	baseURL := fmt.Sprintf("/eu/parl/meps/%d/", m.Identifier)

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, baseURL, "!eurodéputé "+m.Label, "!Fiche descriptive de l'eurodéputée "+m.Label),
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
					render.N("div.headerID", render.Int(m.Identifier)),
				),
				render.N("div.headerTitle", m.Label),
				component.HeaderLangs(translate.Langs, l, ""),
			),

			render.N("main.w"),
			render.N(""),

			component.Footer(l, 0),
		),
	)))
}
