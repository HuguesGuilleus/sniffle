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
		render.N("div.topHeader", tr.PageTop, " (", tr.AboutTextLink, ")"),
		render.N("div.subHeader", render.N("div", "In developpment")),
		render.N("div.titles",
			render.N("div.ids",
				render.N("div.base", "eu/ec/ice"),
				render.N("div.id", ice.Year, "/", ice.Number),
				render.N("div.title", desc.Title)),
		),

		// Translation bar

		render.N("div", "Status: ", ice.Status),
		render.N("div", tr.EU_EC_ICE_ONE.LastUpdate, ice.LastUpdate),
		render.N("div.text", desc.Treaty),

		render.N("div", "Categories: ", strings.Join(ice.Categorie, ", ")),

		render.N("h1", tr.EU_EC_ICE_ONE.H1DescriptionGeneral),
		// website
		render.N("div.text", desc.Objective),
		// image link

		render.N("h1", tr.EU_EC_ICE_ONE.H1DescriptionAnnex),
		render.N("div.text", desc.Annex),

		render.N("h1", tr.EU_EC_ICE_ONE.H1Signature),
		render.N("ul", render.Map(ice.Signature, func(c country.Country, sig uint) render.Node {
			return render.N("li", "Country ", c.String(), ": ", sig)
		}),
		),
	)

	component.Html(t, &page)
}
