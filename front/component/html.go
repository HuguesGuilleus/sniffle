package component

import (
	"sniffle/front"
	"sniffle/front/translate"
	"sniffle/tool/language"
	"sniffle/tool/render"
)

// Meta information about the page
type Page struct {
	Language language.Language

	// <head> informations
	Title       string
	Description string

	// The base URL of the page, without the lang.
	// Ex: /eu/ec/
	BaseURL string
}

var HeadBegin = render.H(`<meta charset=utf-8>` +
	`<meta name=viewport content="width=device-width,initial-scale=1.0">` +
	`<link rel=stylesheet href=/style.` + front.StyleHash + `.css integrity="` + front.StyleIntegrity + `">` +
	`<link rel=icon href=/favicon.ico>`)

// <head> component.
func Head(l language.Language, baseURL, title, description string) render.Node {
	return render.N("head",
		HeadBegin,
		render.N("title", title),
		render.Na("meta", "name", "description").A("content", description),
		langAlternate(baseURL, l, translate.Langs),
	)
}

func langAlternate(baseURL string, pageLang language.Language, langs []language.Language) []render.Node {
	return render.S(langs, "", func(l language.Language) render.Node {
		if pageLang == l {
			return render.Z
		}
		return render.Na("link", "rel", "alternate").
			A("hreflang", l.String()).
			A("href", l.Path(baseURL)).N()
	})
}
