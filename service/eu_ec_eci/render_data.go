package eu_ec_eci

import (
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func renderDataExtraDelay(t *tool.Tool, l language.Language) {
	tr := translate.T[l]
	baseURL := "/eu/ec/eci/data/extradelay/"

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		render.N("head",
			component.HeadBegin,
			render.N("title", tr.EU_EC_ECI.DATA_EXTRADELAY.Name),
			render.Na("meta", "name", "description").A("content", tr.EU_EC_ECI.DATA_EXTRADELAY.Description),
			component.LangAlternate(t.HostURL+baseURL, l, translate.Langs),
		),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerId", "data/extradelay"),
				),
				render.N("div.headerTitle", tr.EU_EC_ECI.DATA_EXTRADELAY.Name),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.w",
				render.N("div.summary", tr.EU_EC_ECI.DATA_EXTRADELAY.Description),
				render.N("ul", render.S(extraDelayData[:], "", func(ice extraDelayICE) render.Node {
					return render.N("li",
						render.N("b",
							render.Na("a", "href", l.Path("../../"+ice.Code+"/")).N(ice.Code),
							" ", ice.Name, render.N("br"),
						),
						" -> ",
						render.S(ice.ExtraDelay, ", ", func(extra ExtraDelay) render.Node {
							return render.Na("a", "href", "https://eur-lex.europa.eu/legal-content/"+l.Upper()+"/TXT/?uri=CELEX:"+extra.CELEX).N(extra.Code)
						}),
					)
				})),
			),
			component.Footer(l, 0),
		),
	)))
}
