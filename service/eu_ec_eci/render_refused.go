package eu_ec_eci

import (
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func renderRefusedIndex(index []*ECIRefused, baseURL string, l language.Language) []byte {
	tr := translate.T[l]
	return render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, baseURL, tr.EU_EC_ECI.REFUSED_INDEX.Name, tr.EU_EC_ECI.REFUSED_INDEX.Description),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerID", "refused"),
				),
				render.N("div.headerTitle", tr.EU_EC_ECI.REFUSED_INDEX.Name),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.w",
				render.N("div.summary",
					render.N("div.edito",
						render.N("div.editoT", tr.GLOBAL.Presentation),
						tr.EU_EC_ECI.REFUSED_INDEX.Description,
					),
					render.Na("a.box", "href", "https://citizens-initiative.europa.eu/find-refused-requests-for-registration_"+l.String()).N(tr.GLOBAL.LinkOfficial),
				),
				component.SearchBlock(l),
				render.S(index, "", func(eci *ECIRefused) render.Node {
					p := ""
					if len(eci.Langs()) != 0 {
						p = eci.Lang.Path(printUint(eci.ID) + "/")
					} else {
						p = eci.OfficielLink()
					}
					return render.Na("a.si.bigItem", "href", p).N(
						render.N("div",
							render.N("div",
								render.N("span.tag.st",
									eci.RefusedDate.In(render.ShortDateZone),
								),
								render.N("span.itemTitle.st", eci.Title),
							),
							render.N("p.itemDesc.st",
								"[", eci.Lang.Human(), "] ",
								eci.PlainDesc,
							),
						),
					)
				}),
			),
			component.Footer(l, component.JsSearch),
		),
	))
}

func renderRefusedOne(eci *ECIRefused) []byte {
	tr := translate.T[eci.Lang]
	ONE := tr.EU_EC_ECI.ONE
	return render.Merge(render.Na("html", "lang", eci.Lang.String()).N(
		render.N("head",
			component.HeadBegin,
			render.N("title", eci.Title),
			render.Na("meta", "name", "description").A("content", eci.PlainDesc),
		),
		render.N("body",
			component.TopHeader(eci.Lang),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(eci.Lang),
					render.N("div.headerID",
						render.Na("a", "href", "..").N("refused"),
						" / ", render.Int(eci.ID)),
				),
				render.N("div.headerTitle", eci.Title),
			),
			render.N("main.wt", component.Toc(eci.Lang), render.N("div.wc",

				render.N("div.summary",
					render.N("div",
						ONE.DescriptionOriginalLangage, eci.Lang.Human(),
					),
					render.N("div",
						render.Na("a.box", "href", eci.OfficielLink()).N(tr.GLOBAL.LinkOfficial),
						render.If(eci.Website != nil, func() render.Node {
							return render.Na("a.box", "href", eci.Website.String()).N(ONE.LinkWebsite)
						}),
					),
					render.If(tool.DevMode, func() render.Node {
						return render.N("div",
							render.N("hr"),
							render.Na("a.box", "href", "./src.json").N("JSON -->"),
						)
					}),
				),

				render.N("h1", tr.EU_EC_ECI.ONE.H1Description),
				render.N("div.text", eci.Objectives),
				render.If(eci.AnnexText != "", func() render.Node {
					return render.N("",
						render.N("h1", tr.EU_EC_ECI.ONE.H1DescriptionAnnex),
						render.N("div.text", eci.AnnexText),
					)
				}),
				renderDoc(eci.Lang, eci.AnnexDoc, tr.EU_EC_ECI.ONE.AnnexDocument),
				renderDoc(eci.Lang, eci.DraftLegal, tr.EU_EC_ECI.ONE.DraftLegal),
				render.If(eci.Treaties != "", func() render.Node {
					return render.N("",
						render.N("h1", tr.EU_EC_ECI.ONE.H1Treaty),
						render.N("p", eci.Treaties),
					)
				}),

				render.N("h1", tr.EU_EC_ECI.ONE.H1Timeline),
				render.N("ol.timeLine",
					render.N("li.timePoint",
						render.N("div.timeHead",
							render.N("span.tag", tr.EU_EC_ECI.Status["REJECTED"]),
							eci.RefusedDate,
						),
						render.Na("a.doc", "href", "https://eur-lex.europa.eu/legal-content/"+eci.Lang.Upper()+"/TXT/?uri=CELEX:"+eci.RefusedCELEX).N(render.N("div.docT", tr.EU_EC_ECI.REFUSED_ONE.RefusalOnline)),
						renderDoc(eci.Lang, &eci.RefusalDocument, tr.EU_EC_ECI.REFUSED_ONE.RefusalDocument),
					),
				),
			)),
			component.Footer(eci.Lang, component.JsToc),
		),
	))
}
