package eu_parl_mep

import (
	"fmt"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func renderTerm(t *tool.Tool, l language.Language, term int, list []mep) {
	// tr := translate.T[l]
	baseURL := fmt.Sprintf("/eu/parl/mep/term-%d/", term)

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, baseURL,
			fmt.Sprintf("!!!Législature %d liste eurodéputés", term),
			fmt.Sprintf("!!!Législature %d liste eurodéputés", term)),

		render.N("body",
			component.InDevHeader(l),
			component.TopHeader(l),

			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerID", "term-", term),
				),
				render.N("div.headerTitle",
					fmt.Sprintf("!!!Législature %d liste eurodéputés", term),
				),
				component.HeaderLangs(translate.Langs, l, ""),
			),

			render.N("main.w",

				// render.Na("a.box", "href", "schema.html").N("schema!!!"),

				render.N("ul",
					render.S(list, "", func(m mep) render.Node {
						return render.N("li", render.Na("a", "href", fmt.Sprintf("../%d/%s.html", m.Identifier, l)).N(
							m.Label, " ",
							render.S(m.Term, " ", func(term int) render.Node {
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
