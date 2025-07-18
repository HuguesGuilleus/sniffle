package eu_parl_mep

import (
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func renderIndex(t *tool.Tool, l language.Language) {
	baseURL := "/eu/parl/mep/"

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, baseURL,
			"!Liste europédutés",
			"!Liste des députées du Parlement européen"),

		render.N("body",
			component.InDevHeader(l),
			component.TopHeader(l),

			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
				),
				render.N("div.headerTitle",
					"!Liste eurodéputés",
				),
				component.HeaderLangs(translate.Langs, l, ""),
			),

			render.N("main.w.home",

				render.N("ul",

					render.N("li", render.Na("a.box", "href", "schema.html").N("schema!!!")),
					render.N("li", render.Na("a.box", "href", "term-0").N("term-0")),
					render.N("li", render.Na("a.box", "href", "term-1").N("term-1")),
					render.N("li", render.Na("a.box", "href", "term-2").N("term-2")),
					render.N("li", render.Na("a.box", "href", "term-3").N("term-3")),
					render.N("li", render.Na("a.box", "href", "term-4").N("term-4")),
					render.N("li", render.Na("a.box", "href", "term-5").N("term-5")),
					render.N("li", render.Na("a.box", "href", "term-6").N("term-6")),
					render.N("li", render.Na("a.box", "href", "term-7").N("term-7")),
					render.N("li", render.Na("a.box", "href", "term-8").N("term-8")),
					render.N("li", render.Na("a.box", "href", "term-9").N("term-9")),
					render.N("li", render.Na("a.box", "href", "term-10").N("term-10")),
				),
			),

			component.Footer(l, 0),
		),
	)))
}
