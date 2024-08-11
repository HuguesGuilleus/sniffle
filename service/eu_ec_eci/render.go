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
			render.N("div",
				render.No("label", render.A("for", "s"), tr.SearchInside),
				render.No("input", render.A("id", "s").A("hidden", "").A("type", "search")),
			),
			render.Slice(eciSlice, func(_ int, eci *ECIOut) render.Node {
				return render.No("a.si.bigItem", render.A("href", fmt.Sprintf("%d/%d/%s.html", eci.Year, eci.Number, l.String())),
					render.N("div",
						render.N("span.tag.st", tr.EU_EC_ECI.Status[eci.Status]),
						render.N("span.box",
							render.N("span.st", eci.Year),
							"/",
							render.N("span.st", eci.Number),
						),
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
				render.N("div", tr.EU_EC_ECI.ONE.Categorie,
					render.Slice(eci.Categorie, func(i int, categorie string) render.Node {
						if i == 0 {
							return render.N("!", tr.EU_EC_ECI.Categorie[categorie])
						}
						return render.N("!", ", "+tr.EU_EC_ECI.Categorie[categorie])
					})),
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
