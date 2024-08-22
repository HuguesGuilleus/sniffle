package component

import (
	"sniffle/front/frontcss"
	"sniffle/tool"
	"sniffle/tool/language"
	"sniffle/tool/render"
)

type Page struct {
	Language    language.Language
	AllLanguage []language.Language

	// <head> informations
	Title       string
	Description string

	// The base URL of the page, without the lang.
	// Ex: /eu/ec/
	BaseURL string

	Body render.Node
}

var htmlHeadBegin = `<meta charset=utf-8>` +
	`<meta name=viewport content="width=device-width,initial-scale=1.0">` +
	`<link rel=stylesheet href=/style.` + frontcss.StyleHash + `.css integrity="` + frontcss.Integrity + `">` +
	`<link rel=icon href=/favicon.ico>`

// Render the page in
func Html(t *tool.Tool, page *Page) {
	data := render.Merge(render.No("html", render.A("lang", page.Language.String()),
		render.N("head",
			render.H(htmlHeadBegin),
			render.N("title", page.Title),
			render.No("meta", render.A("name", "description").A("content", page.Description)),
			langAlternate(t.HostURL+page.BaseURL, page.Language, page.AllLanguage),
		),
		page.Body,
	))
	t.WriteFile(page.BaseURL+page.Language.String()+".html", data)
}

func langAlternate(url string, pageLang language.Language, langs []language.Language) []render.Node {
	return render.Slice(langs, func(_ int, l language.Language) render.Node {
		if pageLang == l {
			return render.Z
		}
		return render.No("link", render.A("rel", "alternate").
			A("hreflang", l.String()).
			A("href", url+l.String()+".html"))
	})
}
