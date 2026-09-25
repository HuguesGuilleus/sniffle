package eu_curia

import (
	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

var schemaHTML = render.Merge(render.Na("html", "lang", "en").N(
	render.N("head",
		component.Head(language.AllEnglish, "/eu/curia/schema.html",
			"EU Curia fetch data",
			"Explain HTTP usage to fetch Curia affairs.",
		),
	),
	render.N("body.edito",
		component.InDevHeader(language.AllEnglish),
		render.N("main.wt.wide",
			component.Toc(language.AllEnglish),
			render.N("div.wc",
				render.N("div.summary", "our usage of the API https://infocuriaws.curia.europa.eu/ to get index and details of European Citizens' Initiative. It is full empiric, be careful!"),

				render.N("h1", "Index"),
				render.N("pre.sch", indexTypes.HTML("")),

				render.N("h1", "Affair"),
				render.N("pre.sch", procedureTypes.HTML("")),
			),
		),
		component.Footer(language.AllEnglish, component.JsSchema|component.JsToc),
	),

	render.N(""),
	render.N(""),
))
