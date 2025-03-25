package eu_ec_eci

import (
	"sniffle/common/country"
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

func renderDataThreshold(t *tool.Tool, l language.Language) {
	baseURL := "/eu/ec/eci/data/threshold/"
	tr := translate.T[l]
	DATA_THRESHOLD := tr.EU_EC_ECI.DATA_THRESHOLD
	countries := translate.Countries(l)

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		render.N("head",
			component.HeadBegin,
			render.N("title", DATA_THRESHOLD.Name),
			render.Na("meta", "name", "description").A("content", DATA_THRESHOLD.Description),
			component.LangAlternate(t.HostURL+baseURL, l, translate.Langs),
		),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerId", "data / threshold"),
				),
				render.N("div.headerTitle", DATA_THRESHOLD.Name),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.wt", component.Toc(l), render.N("div.wc",
				render.N("div.summary", DATA_THRESHOLD.Description),

				render.N("p.noindent", DATA_THRESHOLD.LastCheck, threshold_lastCheck),

				render.S(thresholds[:], "", func(threshold *Threshold) render.Node {
					return render.N("",
						render.N("h1", DATA_THRESHOLD.From, " ", threshold.Begin),
						render.N("div.edito",
							render.N("div.editoT", DATA_THRESHOLD.Calculation),
							tr.EU_EC_ECI.ThresholdRule[threshold.Rule],
						),
						render.N("p.noindent", threshold.Legal.Render(l)),
						render.N("table.right",
							render.N("tr",
								render.N("th", tr.EU_EC_ECI.ONE.Country),
								render.N("th", tr.EU_EC_ECI.ONE.Threshold),
							),
							render.S(countries, "", func(c country.Country) render.Node {
								if threshold.Data[c] == 0 {
									return render.Z
								}
								return render.N("tr",
									render.N("td", tr.Country[c]),
									render.N("td", threshold.Data[c]),
								)
							}),
						),
					)
				}),
			)),
			component.Footer(l, component.JsToc),
		),
	)))
}
