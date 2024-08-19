package eu_ec_eci

import (
	"cmp"
	"encoding/json"
	"fmt"
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
		AllLanguage: t.Languages,
		Title:       tr.EU_EC_ECI.INDEX.Name,
		Description: tr.EU_EC_ECI.INDEX.PageDescription,
		BaseURL:     "/eu/ec/eci/",
	}

	page.Body = render.N("body",
		component.TopHeader(l),
		component.Header(t.Languages, l,
			idNamespace(l),
			render.Z,
			tr.EU_EC_ECI.INDEX.Name),
		render.N("div.w",
			render.N("div.summary",
				render.No("a.box", render.A("href", "https://citizens-initiative.europa.eu/_"+l.String()), tr.LinkOfficial),
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
							A("alt", tr.LogoTitle).A("title", tr.LogoTitle))
					}),
				)
			}),
		),
		component.Footer(l),
	)

	component.Html(t, &page)
}

func renderOne(t *tool.Tool, eci *ECIOut, l language.Language) {
	desc := eci.Description[l]

	page := component.Page{
		Language:    l,
		AllLanguage: t.Languages,
		Title:       desc.Title,
		Description: desc.PlainDesc,
		BaseURL:     fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number),
	}

	tr := translate.AllTranslation[l]
	ONE := tr.EU_EC_ECI.ONE
	page.Body = render.N("body",
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
					render.N("div", ONE.Categorie,
						render.Slice(eci.Categorie, func(i int, categorie string) render.Node {
							if i == 0 {
								return render.N("", tr.EU_EC_ECI.Categorie[categorie])
							}
							return render.N("", ", "+tr.EU_EC_ECI.Categorie[categorie])
						})),
					render.N("div",
						render.No("a.box", render.A("href", fmt.Sprintf(
							"https://citizens-initiative.europa.eu/initiatives/details/%d/%06d_%s", eci.Year, eci.Number, l.String())),
							tr.LinkOfficial),
						render.If(desc.SupportLink != "", func() render.Node {
							return render.No("a.box", render.A("href", desc.SupportLink), ONE.LinkSupport)
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
						A("alt", tr.LogoTitle).A("title", tr.LogoTitle))
				}),

				// Text information
				render.No("h1", render.A("id", "desc"), ONE.H1DescriptionGeneral),
				render.N("h2", ONE.H1DescriptionGeneral),
				render.N("div.text", desc.Objective),

				render.N("h2", ONE.H1DescriptionAnnex),
				render.N("div.text", desc.Annex),

				render.If(desc.Treaty != "", func() render.Node {
					return render.N("",
						render.N("h2", ONE.H1Treaty),
						render.N("div.text", desc.Treaty),
					)
				}),

				// Timeline
				render.N("div.working",
					render.N("h1", ONE.H1Timeline),
					render.If(true, func() render.Node {
						j, _ := json.MarshalIndent(eci.Timeline, "", "\t")
						return render.N("pre", string(j))
					}),
				),

				// Signature
				render.If(len(eci.Signature) > 0, func() render.Node {
					countryOverThreshold := 0
					for c, v := range eci.Signature {
						if eci.Threshold[c] <= v {
							countryOverThreshold++
						}
					}
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
								render.N("span.bigInfoData", countryOverThreshold),
								" / 7",
							),
						),

						render.N("table.right",
							render.N("tr",
								render.N("th", ONE.Country),
								render.N("th", ONE.Signature),
								render.N("th", ONE.Threshold),
							),
							render.If(true, func() render.Node {
								threshold := eci.Threshold
								overtThreshol := render.No("span.tag", render.A("title", ONE.OverThreshold), "✔")

								countrys := make([]country.Country, 0, len(eci.Signature))
								for c := range eci.Signature {
									countrys = append(countrys, c)
								}
								slices.SortFunc(countrys, func(a, b country.Country) int {
									return cmp.Compare(tr.Country[a], tr.Country[b])
								})

								nodes := make([]render.Node, len(countrys))
								for i, c := range countrys {
									sig := eci.Signature[c]
									nodes[i] = render.N("tr",
										render.N("td", tr.Country[c]),
										render.N("td", sig),
										render.N("td",
											render.If(sig >= threshold[c], func() render.Node { return overtThreshol }),
											threshold[c],
										))
								}
								return render.N("", nodes)
							}),
						),
					)
				}),
			),
		),

		component.Footer(l),
	)

	component.Html(t, &page)
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
