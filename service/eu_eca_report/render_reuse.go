package eu_eca_report

import (
	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

var reuseData = render.Merge(render.Na("html", "lang", "en").N(
	render.N("head", component.HeadBegin,
		render.N("title", "Reuse data"),
	),
	render.N("body.data",
		component.TopHeader(language.AllEnglish),
		render.N("header",
			render.N("div.headerID",
				component.HomeAnchor(language.AllEnglish),
				render.Na("a", "href", "/eu/").A("title", translate.T[language.AllEnglish].EU.Name).N("eu"), " / ",
				render.Na("a", "href", "/eu/eca/").A("title", translate.T[language.AllEnglish].EU_ECA.Name).N("eca"),
			),
			render.N("div.headerID", "report / data reuse"),
			render.N("h1", "Reuse data"),
		),
		render.N("p", "Get the index of reports: ", render.Na("a", "href", "/eu/eca/report/index.txt").N("report/index.txt")),
		render.N("p", "Get data: ", render.Na("a", "href", "https://sniffle.eu/eu/eca/report/ID.json").N("https://sniffle.eu/eu/eca/report/ID.json")),
		render.N("p", "Get image: ", render.Na("a", "href", "https://sniffle.eu/eu/eca/report/ID.jpg").N("https://sniffle.eu/eu/eca/report/ID.jpg")),
	),
	component.Footer(language.AllEnglish, component.JsSchema|component.JsToc),
))
