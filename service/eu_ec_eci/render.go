package eu_ec_eci

import (
	"cmp"
	"fmt"
	"slices"
	"sniffle/common/country"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
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

	tr := translate.T[l]
	baseURL := "/eu/ec/eci/"

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, t.HostURL+baseURL, tr.EU_EC_ECI.INDEX.Name, tr.EU_EC_ECI.INDEX.PageDescription),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup", idNamespace(l)),
				render.N("div.headerTitle", tr.EU_EC_ECI.INDEX.Name),
				component.HeaderLangs(l, ""),
			),
			render.N("main.w",
				render.N("div.summary",
					render.Na("a.box", "href", "https://citizens-initiative.europa.eu/find-initiative_"+l.String()).N(tr.LinkOfficial),
					render.Na("a.box", "href", "schema.html").N(tr.SchemaLink),
				),
				render.N("div.searchBlock",
					render.Na("label", "for", "s").N(tr.SearchInside),
					render.Na("input", "id", "s").A("hidden", "").A("type", "search"),
				),
				render.S(eciSlice, "", func(eci *ECIOut) render.Node {
					return render.Na("a.si.bigItem", "href", fmt.Sprintf("%d/%d/%s.html", eci.Year, eci.Number, l)).N(
						render.N("div",
							render.N("span.tag.st", tr.EU_EC_ECI.Status[eci.Status]),
							render.N("span.box.st", render.Int(eci.Year), "/", render.Int(eci.Number)),
						),
						render.N("div.itemTitle.st", eci.Description[l].Title),
						render.N("div", render.S(eci.Categorie, ", ", func(categorie string) render.Node {
							return render.N("span.st", tr.EU_EC_ECI.Categorie[categorie])
						})),
						render.N("p.itemDesc", eci.Description[l].PlainDesc),
						renderImage(eci, true, tr.LogoTitle),
					)
				}),
			),
			component.Footer(l, component.JsSearch),
		),
	)))
}

func renderOne(t *tool.Tool, eci *ECIOut, l language.Language) {
	desc := eci.Description[l]
	tr := translate.T[l]
	ONE := tr.EU_EC_ECI.ONE

	t.WriteFile(l.Path(fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number)), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, fmt.Sprintf("%s/eu/ec/eci/%d/%d/", t.HostURL, eci.Year, eci.Number), desc.Title, desc.PlainDesc),
		render.N("body",
			component.TopHeader(l),
			component.InDevHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerId", render.Int(eci.Year), "/", render.Int(eci.Number)),
				),
				render.N("div.headerTitle", desc.Title),
				component.HeaderLangs(l, ""),
			),
			render.N("main.wt", component.Toc, render.N("div.wc",
				// Summary
				render.N("div.summary",
					render.N("div", ONE.Status, render.N("span.tag", tr.EU_EC_ECI.Status[eci.Status])),
					render.N("div", ONE.LastUpdate, eci.LastUpdate),
					render.N("div", ONE.Categorie, render.S(eci.Categorie, ", ", func(categorie string) render.Node {
						return render.N("", tr.EU_EC_ECI.Categorie[categorie])
					})),
					render.N("div", ONE.DescriptionOriginalLangage, tr.Langage[eci.DescriptionOriginalLangage]),
					render.N("div",
						render.Na("a.box", "href", fmt.Sprintf(
							"https://citizens-initiative.europa.eu/initiatives/details/%d/%06d_%s", eci.Year, eci.Number, l)).
							N(tr.LinkOfficial),
						render.If(desc.FollowUp != nil, func() render.Node {
							return render.Na("a.box", "href", desc.FollowUp.String()).N(ONE.LinkFollowUp)
						}),
						render.If(desc.FollowUp == nil && desc.SupportLink != nil, func() render.Node {
							return render.Na("a.box", "href", desc.SupportLink.String()).N(ONE.LinkSupport)
						}),
						render.If(desc.Website != nil, func() render.Node {
							return render.Na("a.box", "href", desc.Website.String()).N(ONE.LinkWebsite)
						}),
					),
				),

				// Image
				renderImage(eci, false, tr.LogoTitle),

				// Text description
				render.N("h1", ONE.H1Description),
				render.N("div.text", desc.Objective),
				render.IfS(desc.Annex != "" || desc.AnnexDoc != nil, render.N("h2", ONE.H1DescriptionAnnex)),
				render.IfS(desc.Annex != "", render.N("div.text", desc.Annex)),
				desc.AnnexDoc.render(l, ONE.AnnexDocument),
				desc.DraftLegal.render(l, ONE.DraftLegal),
				render.If(desc.Treaty != "", func() render.Node {
					return render.N("",
						render.N("h2", ONE.H1Treaty),
						render.N("p.noindent", desc.Treaty),
					)
				}),

				// Timeline
				render.N("h1", ONE.H1Timeline),
				render.N("ol.timeLine",
					render.S(eci.Timeline, "", func(t Timeline) render.Node {
						child := render.Z
						switch t.Status {
						case "REGISTERED":
							child = t.Register[l].render(l, "Enregistrement")
						case "CLOSED":
							child = render.IfS(t.EarlyClose, render.N("div", ONE.CollectionEarlyClosure))
						case "ANSWERED":
							child = render.N("",
								render.If(t.AnswerResponse != nil, func() render.Node { return t.AnswerResponse[l].render(l, ONE.AnswerKind.Response) }),
								render.If(t.AnswerPressRelease != nil, func() render.Node { return t.AnswerPressRelease[l].render(l, ONE.AnswerKind.PressRelease) }),
								render.If(t.AnswerAnnex != nil, func() render.Node { return t.AnswerAnnex[l].render(l, ONE.AnswerKind.Annex) }),
							)
						case "DEADLINE":
							return render.N("li.timePoint.future", render.N("span.tag", tr.EU_EC_ECI.Status[t.Status]), t.Date)
						}
						return render.N("li.timePoint",
							render.N("div.timeHead", render.N("span.tag", tr.EU_EC_ECI.Status[t.Status]), t.Date),
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

						render.N("div.edito",
							render.N("div.editoT", tr.HELP),
							tr.EU_EC_ECI.ThresholdRule[eci.ThresholdRule],
							" ",
							render.Na("a", "href", "https://citizens-initiative.europa.eu/thresholds_"+l.String()).N(tr.Source),
						),
						render.N("div.bigInfo",
							render.N("div.bifInfoMeta", ONE.CountryOverThreshold),
							render.N("div.bigInfoMain",
								render.N("span.bigInfoData", eci.ThresholdPassTotal),
								" / 7",
							),
						),

						render.N("table.right",
							render.N("tr",
								render.N("th", ONE.Country),
								render.N("th", ONE.Signature),
								render.N("th", ONE.Threshold),
							),
							render.S(eci.countryByName(l), "", func(c country.Country) render.Node {
								return render.N("tr",
									render.N("td", tr.Country[c]),
									render.N("td", eci.Signature[c]),
									render.N("td",
										render.If(eci.ThresholdPass[c], func() render.Node {
											return render.Na("span.tag", "title", ONE.OverThreshold).N("✔")
										}),
										eci.Threshold[c],
									))
							}),
						),
					)
				}),

				// Members
				render.N("div.working", render.N("h1", "$Members ...")),

				// Funding
				render.If(!eci.FundingUpdate.IsZero(), func() render.Node {
					return render.N("",
						render.N("h1", ONE.Funding.Name),
						render.N("div.marginBottom", ONE.LastUpdate, eci.FundingUpdate),
						render.N("div.bigInfo",
							render.N("div.bifInfoMeta", ONE.Funding.Total),
							render.N("div.bigInfoMain.bigInfoData", printEuros(eci.FundingTotal)),
						),
						render.N("table.right",
							render.N("tr",
								render.N("th", ONE.Funding.Sponsor),
								render.N("th", ONE.Funding.Amount),
								render.N("th", ONE.Funding.Date),
							),
							render.S(eci.Sponsor, "", func(s Sponsor) render.Node {
								privateSponsor := render.Na("i", "title", ONE.Funding.PrivateSponsorHelp).N(ONE.Funding.PrivateSponsor)
								return render.N("tr",
									render.N("td", render.IfElse(s.Name != "", func() render.Node {
										return render.N("", s.Name)
									}, func() render.Node { return privateSponsor })),
									render.N("td", printEuros(s.Amount)),
									render.N("td", s.Date.In(render.ShortDateZone)),
								)
							}),
							render.N("caption",
								render.N("div.edito",
									render.N("div.editoT", ONE.Funding.Date),
									ONE.Funding.CaptionDate,
								),
								render.N("div.edito",
									render.N("div.editoT", ONE.Funding.Amount),
									ONE.Funding.CaptionAmount,
								),
							),
						),
						eci.FundingDocument.render(l, ONE.Funding.Document),
					)
				}),
			)),

			component.Footer(l, component.JsToc),
		),
	)))
}

func idNamespace(l language.Language) render.Node {
	tr := translate.T[l]
	return render.N("div.headerId",
		component.HomeAnchor(l),
		render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
		render.Na("a", "href", l.Path("/eu/ec/")).A("title", tr.EU_EC.Name).N("ec"), " / ",
		render.Na("a", "href", l.Path("/eu/ec/eci/")).A("title", tr.EU_EC_ECI.Name).N("eci"),
	)
}

func (doc *Document) render(l language.Language, name render.H) render.Node {
	if doc == nil {
		return render.Z
	}
	tr := translate.T[l]
	return render.Na("a.doc", "href", doc.URL.String()).N(
		render.N("div.docT", name),
		render.If(doc.Language != 0 && doc.Language != l, func() render.Node {
			return render.N("", " ["+tr.Langage[doc.Language]+"]")
		}),
		render.If(doc.Size != 0, func() render.Node {
			return render.N("",
				" (", doc.Size, " ", tr.Byte, ") ", doc.MimeType,
				render.N("div.docName", doc.Name),
			)
		}),
	)
}

func renderImage(eci *ECIOut, needBase bool, title string) render.Node {
	base := "logo"
	if needBase {
		base = fmt.Sprintf("%d/%d/logo", eci.Year, eci.Number)
	}
	return eci.Image.Render(base, title)
}

func printEuros(f float64) any {
	return render.N("",
		int(f),
		fmt.Sprintf(".%02d\u202F€", int64(f*100)%100),
	)
}
