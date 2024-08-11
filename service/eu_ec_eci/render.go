package eu_ec_eci

import (
	"cmp"
	"fmt"
	"slices"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"strings"
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
		component.InDevHeader(l),
		component.Header(t.Languages, l, idNamespace(l),
			render.Z,
			tr.EU_EC_ECI.INDEX.Name),
		render.N("ul.w",
			render.Slice(eciSlice, func(_ int, eci *ECIOut) render.Node {
				return render.N("li",
					render.No("a", render.A("href", fmt.Sprintf("%d/%d/%s.html", eci.Year, eci.Number, l.String())),
						eci.Year, "/", eci.Number,
						render.N("span.tag", tr.EU_EC_ECI.Status[eci.Status]),
						render.IfElse(eci.Description[l] != nil, func() render.Node {
							return render.N("span", eci.Description[l].Title)
						}, func() render.Node {
							return render.N("span", "???")
						}),
						render.N("br"),
						render.If(eci.ImageName != "", func() render.Node {
							return render.No("img.logo", render.
								A("loading", "lazy").
								A("src", fmt.Sprintf("%d/%d/%s", eci.Year, eci.Number, eci.ImageName)).
								A("width", eci.ImageWidth).
								A("height", eci.ImageHeight).
								A("alt", tr.LogoTitle).A("title", tr.LogoTitle))
						}),
					),
				)
			}),
		),
		component.Footer(l),
	)

	component.Html(t, &page)
}

func renderOne(t *tool.Tool, eci *ECIOut, l language.Language) {
	desc := eci.Description[l]

	if desc == nil {
		t.Error("data.err", "err", "eci.Description[] is nil", "lang", l.String(), "year", eci.Year, "nb", eci.Number)
		return
	}

	page := component.Page{
		Language:    l,
		AllLanguage: t.Languages,
		Title:       desc.Title,
		Description: desc.PlainDesc,
		BaseURL:     fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number),
	}

	tr := translate.AllTranslation[l]
	page.Body = render.N("body",
		component.TopHeader(l),
		component.InDevHeader(l),
		component.Header(t.Languages, l, idNamespace(l),
			render.N("div.headerId", eci.Year, "/", eci.Number),
			desc.Title),

		render.N("div.w",
			render.N("div.summary",
				render.N("div", tr.EU_EC_ECI.ONE.Status, render.N("span.tag", tr.EU_EC_ECI.Status[eci.Status])),
				render.N("div", tr.EU_EC_ECI.ONE.LastUpdate, eci.LastUpdate),
				render.N("div", "Categories: ", strings.Join(eci.Categorie, ", ")),
				render.N("div",
					render.No("a.box", render.A("href", fmt.Sprintf(
						"https://citizens-initiative.europa.eu/initiatives/details/%d/%06d_%s", eci.Year, eci.Number, l.String())),
						tr.EU_EC_ECI.ONE.LinkOfficial),
					render.If(desc.SupportLink != "", func() render.Node {
						return render.No("a.box", render.A("href", desc.SupportLink), tr.EU_EC_ECI.ONE.LinkSupport)
					}),
					render.If(desc.Website != nil, func() render.Node {
						return render.No("a.box", render.A("href", desc.Website.String()), tr.EU_EC_ECI.ONE.LinkWebsite)
					}),
				),
			),

			render.If(eci.ImageName != "", func() render.Node {
				return render.No("img.logo.marginTop", render.
					A("loading", "lazy").
					A("src", eci.ImageName).
					A("width", eci.ImageWidth).
					A("height", eci.ImageHeight).
					A("alt", tr.LogoTitle).A("title", tr.LogoTitle))
			}),

			render.N("h1", tr.EU_EC_ECI.ONE.H1DescriptionGeneral),
			render.N("div.text", desc.Objective),
			// image link

			render.N("h1", tr.EU_EC_ECI.ONE.H1DescriptionAnnex),
			render.N("div.text", desc.Annex),

			render.N("h1", tr.EU_EC_ECI.ONE.H1Treaty),
			render.N("div.text", desc.Treaty),

			render.N("h1", tr.EU_EC_ECI.ONE.H1Signature),
			render.N("ul", render.Map(eci.Signature, func(c country.Country, sig uint) render.Node {
				return render.N("li", "Country ", c.String(), ": ", sig)
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
