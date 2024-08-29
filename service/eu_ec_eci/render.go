package eu_ec_eci

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/render"
)

func renderIndex(t *tool.Tool, eciSlice []*ECIOut, l language.Language) {
	slices.SortFunc(eciSlice, func(a, b *ECIOut) int {
		return cmp.Or(
			cmp.Compare(b.Year, a.Year),
			cmp.Compare(b.Number, a.Number),
		)
	})

	tr := translate.AllTranslation[l]

	page := component.Page{
		Language:    l,
		Title:       tr.EU_EC_ECI.INDEX.Name,
		Description: tr.EU_EC_ECI.INDEX.PageDescription,
		BaseURL:     "/eu/ec/eci/",
	}

	component.HTML(t, &page, render.N("body",
		component.TopHeader(l),
		component.Header(t.Languages, l,
			idNamespace(l),
			render.Z,
			tr.EU_EC_ECI.INDEX.Name),
		render.N("div.w",
			render.N("div.summary",
				render.No("a.box", render.A("href", "https://citizens-initiative.europa.eu/find-initiative_"+l.String()), tr.LinkOfficial),
			),
			render.N("div.searchBlock",
				render.No("label", render.A("for", "s"), tr.SearchInside),
				render.No("input", render.A("id", "s").A("hidden", "").A("type", "search")),
			),
			render.Slice(eciSlice, func(_ int, eci *ECIOut) render.Node {
				return render.No("a.si.bigItem", render.A("href", fmt.Sprintf("%d/%d/%s.html", eci.Year, eci.Number, l.String())),
					render.N("div",
						render.N("span.tag.st", tr.EU_EC_ECI.Status[eci.Status]),
						render.N("span.box.st", render.Int(eci.Year), "/", render.Int(eci.Number)),
					),
					render.N("div.itemTitle.st", eci.Description[l].Title),
					render.N("div", render.SliceSeparator(eci.Categorie, ", ", func(_ int, categorie string) render.Node {
						return render.N("span.st", tr.EU_EC_ECI.Categorie[categorie])
					})),
					render.N("p.itemDesc", eci.Description[l].PlainDesc),
					render.If(eci.ImageName != "", func() render.Node {
						return render.No("img.logo", render.
							A("loading", "lazy").
							A("src", fmt.Sprintf("%d/%d/%s", eci.Year, eci.Number, eci.ImageName)).
							A("width", eci.ImageWidth).
							A("height", eci.ImageHeight).
							A("title", tr.LogoTitle))
					}),
				)
			}),
		),
		component.Footer(l, component.JsSearch),
	))
}

func renderOne(t *tool.Tool, eci *ECIOut, l language.Language) {
	desc := eci.Description[l]
	tr := translate.AllTranslation[l]
	ONE := tr.EU_EC_ECI.ONE
	countrysByName := slices.SortedFunc(maps.Keys(eci.Signature), func(a, b country.Country) int {
		return cmp.Compare(tr.Country[a], tr.Country[b])
	})

	page := component.Page{
		Language:    l,
		Title:       desc.Title,
		Description: desc.PlainDesc,
		BaseURL:     fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number),
	}

	component.HTML(t, &page, render.N("body",
		component.TopHeader(l),
		component.InDevHeader(l),
		component.Header(t.Languages, l, idNamespace(l),
			render.N("div.headerId", render.Int(eci.Year), "/", render.Int(eci.Number)),
			desc.Title),

		render.N("div.wt",
			render.N("div", render.No("div", render.A("id", "toc"))),

			render.N("div.wc",
				// Summary
				render.N("div.summary",
					render.N("div", ONE.Status, render.N("span.tag", tr.EU_EC_ECI.Status[eci.Status])),
					render.N("div", ONE.LastUpdate, eci.LastUpdate),
					render.N("div", ONE.Categorie, render.SliceSeparator(eci.Categorie, ", ", func(_ int, categorie string) render.Node {
						return render.N("", tr.EU_EC_ECI.Categorie[categorie])
					})),
					render.N("div", ONE.DescriptionOriginalLangage, tr.Langage[eci.DescriptionOriginalLangage]),
					render.N("div",
						render.No("a.box", render.A("href", fmt.Sprintf(
							"https://citizens-initiative.europa.eu/initiatives/details/%d/%06d_%s", eci.Year, eci.Number, l.String())),
							tr.LinkOfficial),
						render.If(desc.SupportLink != nil, func() render.Node {
							return render.No("a.box", render.A("href", desc.SupportLink.String()), ONE.LinkSupport)
						}),
						render.If(desc.Website != nil, func() render.Node {
							return render.No("a.box", render.A("href", desc.Website.String()), ONE.LinkWebsite)
						}),
					),
				),

				// Image
				render.If(eci.ImageName != "", func() render.Node {
					return render.No("img.logo.marginTop", render.
						A("loading", "lazy").
						A("src", eci.ImageName).
						A("width", eci.ImageWidth).
						A("height", eci.ImageHeight).
						A("title", tr.LogoTitle))
				}),

				// Text information
				render.N("h1", ONE.H1Description),
				render.N("div.text", desc.Objective),
				render.IfS(desc.Annex != "" || desc.AnnexDoc != nil, render.N("h2", ONE.H1DescriptionAnnex)),
				render.IfS(desc.Annex != "", render.N("div.text", desc.Annex)),
				render.If(desc.AnnexDoc != nil, func() render.Node {
					return desc.AnnexDoc.render(l, ONE.AnnexDocument)
				}),
				render.If(desc.Treaty != "", func() render.Node {
					return render.N("",
						render.N("h2", ONE.H1Treaty),
						render.N("p.noindent", desc.Treaty),
					)
				}),

				// Timeline
				render.N("h1", ONE.H1Timeline),
				render.N("ol.timeLine",
					render.Slice(eci.Timeline, func(_ int, t Timeline) render.Node {
						child := render.Z
						switch t.Status {
						case "REGISTERED":
							child = t.Register[l].render(l, "Enregistrement")
						case "CLOSED":
							child = render.IfS(t.EarlyClose, render.N("div", ONE.CollectionEarlyClosure))
						case "ANSWERED":
							// TODO: answer document
						case "DEADLINE":
							return render.N("li.timePoint.future", render.N("span.tag", tr.EU_EC_ECI.Status[t.Status]), t.Date)
						}
						return render.N("li.timePoint",
							render.N("span.tag", tr.EU_EC_ECI.Status[t.Status]), t.Date,
							child)
					}),
				),

				// Signature
				render.If(len(eci.Signature) != 0, func() render.Node {
					return render.N("",
						render.N("h1", tr.EU_EC_ECI.ONE.H1Signature),
						render.If(eci.ValidatedSignature, func() render.Node {
							return render.N("div.marginBottom", ONE.ValidatedSignature)
						}),
						render.If(!eci.PaperSignaturesUpdate.IsZero(), func() render.Node {
							return render.N("div.marginBottom", ONE.PaperSignaturesUpdate, eci.PaperSignaturesUpdate)
						}),
						render.N("div.bigInfo",
							render.N("div.bifInfoMeta", ONE.SignatureSum),
							render.N("div.bigInfoMain",
								render.N("span.bigInfoData", eci.TotalSignature),
								" / 1 000 000",
							),
						),
						// render.N("div.edito",
						// 	render.N("div.editoT", tr.HELP),
						// 	render.N("p", "Les seuils correspondent au nombre de députés au Parlement européen élus dans chaque État membre, multiplié par le nombre total de députés au Parlement européen.")),
						render.N("div.bigInfo",
							render.N("div.bifInfoMeta", ONE.CountryOverThreshold),
							render.N("div.bigInfoMain",
								render.N("span.bigInfoData", eci.ThresholdPassed),
								" / 7",
							),
						),

						render.N("table.right",
							render.N("tr",
								render.N("th", ONE.Country),
								render.N("th", ONE.Signature),
								render.N("th", ONE.Threshold),
							),
							render.Slice(countrysByName, func(_ int, c country.Country) render.Node {
								sig := eci.Signature[c]
								return render.N("tr",
									render.N("td", tr.Country[c]),
									render.N("td", sig),
									render.N("td",
										render.If(sig >= eci.Threshold[c], func() render.Node {
											return render.No("span.tag", render.A("title", ONE.OverThreshold), "✔")
										}),
										eci.Threshold[c],
									))
							}),
						),
					)
				}),
			),
		),

		component.Footer(l, component.JsToc),
	))
}

func idNamespace(l language.Language) render.Node {
	tr := translate.AllTranslation[l]
	return render.N("div.headerId",
		component.HomeAnchor(l),
		render.No("a",
			render.A("href", "/eu/"+l.String()+".html").A("title", tr.EU.Name),
			"eu"),
		" / ",
		render.No("a",
			render.A("href", "/eu/ec/"+l.String()+".html").A("title", tr.EU_EC.Name),
			"ec"),
		" / ",
		render.No("a",
			render.A("href", "/eu/ec/eci/"+l.String()+".html").A("title", tr.EU_EC_ECI.Name),
			"eci"),
	)
}

func (doc *Document) render(lang language.Language, name render.H) render.Node {
	return render.No("a.doc", render.A("href", doc.URL.String()),
		render.N("div.docT", name, render.H(" &gt;&gt;&gt;")),
		render.If(doc.Size != 0, func() render.Node {
			tr := translate.AllTranslation[lang]
			return render.N("",
				tr.Langage[doc.Language], " (", doc.Size, " ", tr.Byte, ") ", doc.MimeType,
				render.N("div.docName", doc.Name))
		}),
	)
}
