package component

import (
	"sniffle/common/language"
	"sniffle/front/translate"
	"sniffle/tool/render"
)

func TopHeader(l language.Language) render.Node {
	return render.N("div.topHeader",
		translate.T[l].GLOBAL.NotEUWebsite,
		render.H(" ("),
		render.Na("a", "href", l.Path("/about/")).N(translate.T[l].GLOBAL.AboutTextLink),
		render.H(")"))
}

// A header to indicated taht this page is currently in development.
func InDevHeader(l language.Language) render.Node {
	return render.N("div.subHeader", translate.T[l].GLOBAL.InDev)
}

func HomeAnchor(l language.Language) render.Node {
	tr := translate.T[l]
	return render.N("",
		render.Na("a.headerHome", "href", l.Path("/")).A("title", tr.HOME.Name).N("⾕"),
		" / ",
	)
}

// Links to different
func HeaderLangs(langs []language.Language, pageLang language.Language, basePath string) render.Node {
	return render.N("div.headerLangs", render.S(langs, "", func(l language.Language) render.Node {
		if pageLang == l {
			return render.Na("span.headerOneLang", "title", l.Human()).N(l.String())
		}
		return render.Na("a.headerOneLang", "title", l.Human()).
			A("hreflang", l.String()).
			A("href", l.Path(basePath)).
			N(l.String())
	}))
}

// The table of content node.
// Yout need JS to fill it.
func Toc(l language.Language) render.Node {
	return render.N("div", render.Na("div", "id", "toc").N(
		render.Na("a.wi", "href", "#").N(translate.T[l].GLOBAL.PageTop),
	))
}

func SearchBlock(l language.Language) render.Node {
	return render.Na("div.searchBlock", "hidden", "").N(
		render.Na("label", "for", "search").N(translate.T[l].GLOBAL.SearchInside),
		render.Na("input", "id", "search").A("type", "search"),
	)
}
