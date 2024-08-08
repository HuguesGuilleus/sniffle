package eu_ec_ice

import (
	"fmt"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"strings"
)

func renderOne(t *tool.Tool, ice *ICEOut, l language.Language) {
	desc := ice.Description[l]

	if desc == nil {
		t.Error("data.err", "err", "ice.Description[] is nil", "lang", l.String(), "year", ice.Year, "nb", ice.Number)
		return
	}

	page := component.Page{
		Language:    l,
		AllLanguage: t.Languages,
		Title:       desc.Title,
		HostURL:     t.HostURL,
		BaseURL:     fmt.Sprintf("/eu/ec/ice/%d/%d/", ice.Year, ice.Number),
	}

	tr := translate.AllTranslation[l]
	page.Body = render.N("body",
		component.TopHeader(l),
		component.InDevHeader(l),
		render.N("header",
			render.N("div.headerSup",
				render.N("div.headerId",
					render.No("a",
						render.A("href", "/").A("title", tr.HomeTitle),
						tr.HomeName),
					" / ",
					render.No("a",
						render.A("href", "/eu/").A("title", tr.EU.Name),
						"eu"),
					" / ",
					render.No("a",
						render.A("href", "/eu/ec/").A("title", tr.EU_EC.Name),
						"ec"),
					" / ",
					render.No("a",
						render.A("href", "/eu/ec/ice/").A("title", tr.EU_EC_ICE.Name),
						"ice"),
				),
				render.N("div.headerId", ice.Year, "/", ice.Number)),
			render.N("div.headerTitle", desc.Title),
		),

		// Translation bar

		render.N("div.w",
			render.N("div.summary",
				render.N("div", "Status: ", render.N("span.tag", ice.Status)),
				render.N("div", tr.EU_EC_ICE.ONE.LastUpdate, ice.LastUpdate),
				render.N("div", "Categories: ", strings.Join(ice.Categorie, ", ")),
			),

			render.N("h1", tr.EU_EC_ICE.ONE.H1DescriptionGeneral),
			// website
			render.N("div.text", desc.Objective),
			// image link

			render.N("h1", tr.EU_EC_ICE.ONE.H1DescriptionAnnex),
			render.N("div.text", desc.Annex),

			render.N("h1", tr.EU_EC_ICE.ONE.H1Treaty),
			render.N("div.text", desc.Treaty),

			render.N("h1", tr.EU_EC_ICE.ONE.H1Signature),
			render.N("ul", render.Map(ice.Signature, func(c country.Country, sig uint) render.Node {
				return render.N("li", "Country ", c.String(), ": ", sig)
			}),
			),
		),
	)

	component.Html(t, &page)
}
