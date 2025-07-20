package eu_ec_eci

import (
	"fmt"
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
	"strconv"
)

func writeIndex(t *tool.Tool, eciByYear map[uint][]*ECIOut, l language.Language) {
	tr := translate.T[l]
	baseURL := "/eu/ec/eci/"

	t.WriteFile(l.Path(baseURL), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, baseURL, tr.EU_EC_ECI.INDEX.Name, tr.EU_EC_ECI.INDEX.PageDescription),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup", idNamespace(l)),
				render.N("div.headerTitle", tr.EU_EC_ECI.INDEX.Name),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.wt", component.Toc(l), render.N("div.wc",
				render.N("div.summary",
					render.N("div.edito",
						render.N("div.editoT", tr.GLOBAL.Presentation),
						tr.EU_EC_ECI.INDEX.Help,
					),
					render.N("div.boxFlex",
						render.Na("a.box", "href", "https://citizens-initiative.europa.eu/_"+l.String()).N(tr.GLOBAL.LinkOfficial),
						render.Na("a.box", "href", "https://citizens-initiative.europa.eu/find-initiative_"+l.String()).N(tr.EU_EC_ECI.INDEX.IndexLink),
					),
					render.N("hr"),
					render.N("div.boxFlex",
						render.Na("a.box", "href", l.Path("refused/")).N(tr.EU_EC_ECI.REFUSED_INDEX.Name),
					),
					render.N("hr"),
					render.N("div.boxFlex",
						render.Na("a.box", "href", "schema.html").N(tr.GLOBAL.SchemaLink),
						render.Na("a.box", "href", l.Path("data-extradelay/")).N(tr.EU_EC_ECI.DATA_EXTRADELAY.Name),
						render.Na("a.box", "href", l.Path("data-threshold/")).N(tr.EU_EC_ECI.DATA_THRESHOLD.Name),
					),
				),
				component.SearchBlock(l),
				render.MapReverse(eciByYear, func(year uint, slice []*ECIOut) render.Node {
					return render.N("div.sg",
						render.N("h1", render.Int(year)),
						render.S(slice, "", func(eci *ECIOut) render.Node {
							l := l
							if eci.Description[l] == nil {
								l = eci.OriginalLangage
							}
							return render.Na("a.si.bigItem", "href", fmt.Sprintf("%d/%d/%s.html", eci.Year, eci.Number, l)).N(
								render.N("div",
									render.N("span.tag.st", render.Int(eci.Year), "/", render.Int(eci.Number)),
									render.N("span.tag.st", tr.EU_EC_ECI.Status[eci.Status]),
								),
								render.N("div.itemTitle.st", eci.Description[l].Title),
								render.N("div", render.S(eci.Categorie, ", ", func(categorie string) render.Node {
									return render.N("span.st", tr.EU_EC_ECI.Categorie[categorie])
								})),
								render.N("p.itemDesc", eci.Description[l].PlainDesc),
								renderImage(eci, true, tr.GLOBAL.LogoTitle),
							)
						}),
					)
				}),
			)),
			component.Footer(l, component.JsSearch|component.JsToc),
		),
	)))
}

func writeOne(t *tool.Tool, eci *ECIOut, l language.Language) {
	desc := eci.Description[l]
	tr := translate.T[l]
	ONE := tr.EU_EC_ECI.ONE

	t.WriteFile(l.Path(fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number)), render.Merge(render.Na("html", "lang", l.String()).N(
		render.N("head",
			component.HeadBegin,
			render.N("title", desc.Title),
			render.Na("meta", "name", "description").A("content", desc.PlainDesc),
			component.LangAlternate(fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number), l, eci.Langs()),
		),
		render.N("body",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerID", render.Int(eci.Year), "/", render.Int(eci.Number)),
				),
				render.N("div.headerTitle", desc.Title),
				component.HeaderLangs(eci.Langs(), l, ""),
			),
			render.N("main.wt", component.Toc(l), render.N("div.wc",
				// Summary
				render.N("div.summary",
					render.N("div", ONE.Status, render.N("span.tag", tr.EU_EC_ECI.Status[eci.Status])),
					render.N("div", ONE.LastUpdate, eci.LastUpdate),
					render.N("div", ONE.Categorie, render.S(eci.Categorie, ", ", func(categorie string) render.Node {
						return render.N("", tr.EU_EC_ECI.Categorie[categorie])
					})),
					render.N("div", ONE.DescriptionOriginalLangage, tr.Langage[eci.OriginalLangage]),
					render.N("div.boxFlex",
						render.Na("a.box", "href", fmt.Sprintf(
							"https://citizens-initiative.europa.eu/initiatives/details/%d/%06d_%s", eci.Year, eci.Number, l)).
							N(tr.GLOBAL.LinkOfficial),
						render.If(desc.SupportLink != nil, func() render.Node {
							return render.Na("a.box", "href", desc.SupportLink.String()).N(ONE.LinkSignature)
						}),
						render.If(desc.FollowUp != nil, func() render.Node {
							return render.Na("a.box", "href", desc.FollowUp.String()).N(ONE.LinkFollowUp)
						}),
						render.If(desc.Website != nil, func() render.Node {
							return render.Na("a.box", "href", desc.Website.String()).N(ONE.LinkWebsite)
						}),
					),
					render.If(tool.DevMode, func() render.Node {
						return render.N("div",
							render.N("hr"),
							render.Na("a.box", "href", "./src.json").N("JSON -->"),
						)
					}),
				),

				// Image
				renderImage(eci, false, tr.GLOBAL.LogoTitle),

				// Text description
				render.N("h1", ONE.H1Description),
				render.N("div.text", desc.Objective),
				render.If(desc.Annex != "", func() render.Node {
					return render.N("",
						render.N("h2", ONE.H1DescriptionAnnex),
						render.N("div.text", desc.Annex),
					)
				}),
				renderDoc(l, desc.AnnexDoc, ONE.AnnexDocument),
				renderDoc(l, desc.DraftLegal, ONE.DraftLegal),
				render.If(desc.Treaty != "", func() render.Node {
					return render.N("",
						render.N("h2", ONE.H1Treaty),
						render.N("p.noindent", desc.Treaty),
					)
				}),

				// Timeline
				render.N("h1", ONE.H1Timeline),
				render.N("ol.timeLine",
					render.S(eci.Timeline, "", func(e Event) render.Node {
						child := render.Z
						switch e.Status {
						case "REGISTERED":
							child = renderDoc(l, e.Register[l], ONE.Registration)
							if e.RegisterCorrigendum[l] != nil {
								child = render.N("", child, renderDoc(l, e.RegisterCorrigendum[l], ONE.RegistrationCorrigendum))
							}
						case "CLOSED":
							if e.EarlyClose {
								child = render.N("div", ONE.CollectionEarlyClosure)
							}
							if len(e.ExtraDelay) != 0 {
								child = render.N("", child, render.N("div",
									ONE.ExtraDelay,
									render.S(e.ExtraDelay, ", ", func(extra component.Legal) render.Node {
										return extra.Render(l)
									}),
								))
							}
						case "ANSWERED":
							child = render.N("",
								render.If(e.AnswerResponse != nil, func() render.Node {
									return renderDoc(l, e.AnswerResponse[l], ONE.AnswerKind.Response)
								}),
								render.If(e.AnswerPressRelease != nil, func() render.Node {
									return renderDoc(l, e.AnswerPressRelease[l], ONE.AnswerKind.PressRelease)
								}),
								render.If(e.AnswerAnnex != nil, func() render.Node {
									return renderDoc(l, e.AnswerAnnex[l], ONE.AnswerKind.Annex)
								}),
							)
						case "DEADLINE":
							return render.N("li.timePoint.future", render.N("span.tag", tr.EU_EC_ECI.Status[e.Status]), e.Date)
						}
						return render.N("li.timePoint",
							render.N("div.timeHead", render.N("span.tag", tr.EU_EC_ECI.Status[e.Status]), e.Date),
							child)
					}),
				),

				// Signature
				render.If(len(eci.Signature) != 0, func() render.Node {
					return render.N("",
						render.N("h1", tr.EU_EC_ECI.ONE.H1Signature),
						render.If(!eci.PaperSignaturesUpdate.IsZero(), func() render.Node {
							return render.N("div.marginBottom", ONE.PaperSignaturesUpdate, eci.PaperSignaturesUpdate)
						}),

						render.N("div.bigInfo",
							render.N("div.bigInfoMeta", ONE.SignatureSum),
							render.N("div.bigInfoMain",
								render.N("span.bigInfoData", eci.TotalSignature),
								" / 1 000 000",
							),
						),

						render.N("div.bigInfo",
							render.N("div.bigInfoMeta", ONE.CountryOverThreshold),
							render.N("div.bigInfoMain",
								render.N("span.bigInfoData", eci.ThresholdPassTotal),
								" / 7",
							),
							render.N("div.edito",
								render.N("div.editoT", tr.GLOBAL.HELP),
								render.N("div.marginBottom", tr.EU_EC_ECI.ThresholdRule[eci.Threshold.Rule]),
								tr.GLOBAL.Source, ": ",
								eci.Threshold.Legal.Render(l),
							),
						),

						render.N("table.right",
							render.N("tr",
								render.N("th", ONE.Country),
								render.N("th", ONE.Signature),
								render.N("th", ONE.Threshold),
							),
							render.S(eci.Signature, "", func(sig Signature) render.Node {
								afterDelayed := render.Na("span.tag", "title", ONE.AfterSubmission).N("✘")
								thresholdPass := render.Na("span.tag", "title", ONE.OverThreshold).N("✔")
								return render.N("tr",
									render.N("td", tr.Country[sig.Country]),
									render.N("td",
										render.IfS(sig.After, afterDelayed),
										sig.Count,
									),
									render.N("td",
										render.IfS(sig.ThresholdPass, thresholdPass),
										sig.Threshold,
									),
								)
							}),
						),
					)
				}),

				// Members
				render.N("h1", ONE.Member.H1),
				render.N("ul.peopleIndex", render.S(eci.Members, "", func(m Member) render.Node {
					return renderMember(l, &m)
				})),

				// Funding
				render.If(!eci.FundingUpdate.IsZero(), func() render.Node {
					return render.N("",
						render.N("h1", ONE.Funding.Name),
						render.N("div.marginBottom", ONE.LastUpdate, eci.FundingUpdate),
						render.N("div.bigInfo",
							render.N("div.bigInfoMeta", ONE.Funding.Total),
							render.N("div.bigInfoMain.bigInfoData", printEuros(eci.FundingTotal)),
						),
						render.N("table.right",
							render.N("tr",
								render.N("th", ONE.Funding.Sponsor),
								render.N("th", ONE.Funding.Kind),
								render.N("th", ONE.Funding.Amount, "*"),
								render.N("th", ONE.Funding.Date, "*"),
							),
							render.S(eci.Sponsor, "", func(s Sponsor) render.Node {
								anonymous := render.Na("i", "title", ONE.Funding.AnonymousHelp).N(ONE.Funding.Anonymous)
								return render.N("tr",
									render.N("td", render.IfElse(s.Name != "",
										func() render.Node { return render.N("", s.Name) },
										func() render.Node { return anonymous },
									)),
									render.N("td", render.IfElse(s.IsPrivate,
										func() render.Node { return render.N("", ONE.Funding.KindPrivate) },
										func() render.Node { return render.N("", ONE.Funding.KindOrganisation) },
									)),
									render.N("td", printEuros(s.Amount)),
									render.N("td", s.Date.In(render.ShortDateZone)),
								)
							}),
							render.N("caption",
								render.N("div.edito",
									render.N("div.editoT", ONE.Funding.Amount),
									ONE.Funding.CaptionAmount,
								),
								render.N("div.edito",
									render.N("div.editoT", ONE.Funding.Date),
									ONE.Funding.CaptionDate,
								),
							),
						),
						renderDoc(l, eci.FundingDocument, ONE.Funding.Document),
					)
				}),
			)),

			component.Footer(l, component.JsToc),
		),
	)))
}

func idNamespace(l language.Language) render.Node {
	tr := translate.T[l]
	return render.N("div.headerID",
		component.HomeAnchor(l),
		render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
		render.Na("a", "href", l.Path("/eu/ec/")).A("title", tr.EU_EC.Name).N("ec"), " / ",
		render.Na("a", "href", l.Path("/eu/ec/eci/")).A("title", tr.EU_EC_ECI.Name).N("eci"),
	)
}

func renderDoc(l language.Language, doc *Document, name render.H) render.Node {
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
				" (", doc.Size, " ", tr.GLOBAL.Byte, ") ", doc.MimeType,
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

func renderMember(l language.Language, m *Member) render.Node {
	ONE := translate.T[l].EU_EC_ECI.ONE
	return render.N("li.people",
		render.N("div",
			render.N("span.tag", ONE.Member.Type[m.Type]),
			render.N("em", m.FullName),
		),
		render.If(m.HrefURL != "", func() render.Node {
			return render.N("div", render.Na("a", "href", m.HrefURL).N(m.DisplayURL))
		}),
		render.N("div",
			render.If(!m.Start.IsZero(), func() render.Node {
				return render.N("", ONE.Member.Start, m.Start, ". ")
			}),
			render.If(!m.End.IsZero(), func() render.Node {
				return render.N("", ONE.Member.End, m.End, ". ")
			}),
			render.If(m.ResidenceCountry.NotZero(), func() render.Node {
				return render.N("",
					ONE.Member.Country,
					translate.T[l].Country[m.ResidenceCountry],
				)
			}),
		),
		render.If(m.Replaced != nil, func() render.Node {
			return render.N("ul.peopleIndex", renderMember(l, m.Replaced))
		}),
	)
}

func printEuros(f float64) any {
	return render.N("",
		int(f),
		fmt.Sprintf(".%02d\u202F€", int64(f*100)%100),
	)
}

func printUint(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}
