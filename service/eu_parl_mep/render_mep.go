package eu_parl_mep

import (
	"fmt"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func renderMeps(t *tool.Tool, l language.Language, m mep) {
	// tr := translate.T[l]
	baseURL := fmt.Sprintf("/eu/parl/mep/%d/", m.Identifier)

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, baseURL, "!eurodéputé "+m.Label, "!Fiche descriptive de l'eurodéputée "+m.Label),
		render.N("body",
			component.InDevHeader(l),
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
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
