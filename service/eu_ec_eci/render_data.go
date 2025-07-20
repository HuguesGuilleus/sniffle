package eu_ec_eci

import (
	"sniffle/common/country"
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func renderDataExtraDelay(t *tool.Tool, eciByYear map[uint][]*ECIOut, l language.Language) {
	titles := make(map[uint]string)
	for _, eci := range extraDelayData {
		titles[eci.ID] = eci.Name
	}
	for _, s := range eciByYear {
		for _, eci := range s {
			if desc := eci.Description[l]; desc != nil {
				titles[eci.ID] = desc.Title
			}
		}
	}

	tr := translate.T[l]
	baseURL := "/eu/ec/eci/data-extradelay/"

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		render.N("head",
			component.HeadBegin,
			render.N("title", tr.EU_EC_ECI.DATA_EXTRADELAY.Name),
			render.Na("meta", "name", "description").A("content", tr.EU_EC_ECI.DATA_EXTRADELAY.Description),
			component.LangAlternate(baseURL, l, translate.Langs),
		),
		render.N("body.edito",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerID", "data/extradelay"),
				),
				render.N("div.headerTitle", tr.EU_EC_ECI.DATA_EXTRADELAY.Name),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.w",
				render.N("div.summary", tr.EU_EC_ECI.DATA_EXTRADELAY.Description),
				render.N("ul", render.S(extraDelayData[:], "", func(eci extraDelayICE) render.Node {
					return render.N("li",
						render.N("b",
							render.Na("a", "href", l.Path("../../"+eci.Code+"/")).N(eci.Code),
							" ", titles[eci.ID], render.N("br"),
						),
						" -> ",
						render.S(eci.ExtraDelay, ", ", func(extra component.Legal) render.Node {
							return extra.Render(l)
						}),
					)
				})),
			),
			component.Footer(l, 0),
		),
	)))
}

func renderDataThreshold(t *tool.Tool, l language.Language) {
	baseURL := "/eu/ec/eci/data-threshold/"
	tr := translate.T[l]
	DATA_THRESHOLD := tr.EU_EC_ECI.DATA_THRESHOLD
	countries := translate.Countries(l)

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		render.N("head",
			component.HeadBegin,
			render.N("title", DATA_THRESHOLD.Name),
			render.Na("meta", "name", "description").A("content", DATA_THRESHOLD.Description),
			component.LangAlternate(baseURL, l, translate.Langs),
		),
		render.N("body.edito",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerID", "data / threshold"),
				),
				render.N("div.headerTitle", DATA_THRESHOLD.Name),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.wt.wide", component.Toc(l), render.N("div.wc",
				render.N("div.summary", DATA_THRESHOLD.Description),

				render.N("p.noindent", DATA_THRESHOLD.LastCheck, threshold_lastCheck),

				render.N("h1", DATA_THRESHOLD.H1Data),
				render.N("table.right",
					render.N("tr",
						render.N("th", tr.EU_EC_ECI.ONE.Country),
						render.S(thresholds[:], "", func(t *Threshold) render.Node {
							return render.N("th", DATA_THRESHOLD.From, " ", t.Begin)
						}),
					),
					render.S(countries, "", func(c country.Country) render.Node {
						return render.N("tr",
							render.N("td", tr.Country[c]),
							render.S(thresholds[:], "", func(t *Threshold) render.Node {
								if t.Data[c] == 0 {
									return render.N("td.blanck")
								}
								return render.N("td", t.Data[c])
							}),
						)
					}),
				),

				render.N("div.subw",
					render.N("h1", DATA_THRESHOLD.H1Rule),
					render.S(thresholds[:], "", func(threshold *Threshold) render.Node {
						return render.N("",
							render.N("h2", DATA_THRESHOLD.From, " ", threshold.Begin),
							render.N("div.edito",
								render.N("div.editoT", DATA_THRESHOLD.Calculation),
								tr.EU_EC_ECI.ThresholdRule[threshold.Rule],
							),
							render.N("p.noindent", threshold.Legal.Render(l)),
						)
					}),
				),
			)),
			component.Footer(l, component.JsToc),
		),
	)))
}
