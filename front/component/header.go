package component

import (
	"sniffle/front/translate"
	"sniffle/tool/language"
	"sniffle/tool/render"
)

func TopHeader(l language.Language) render.Node {
	return render.N("div.topHeader",
		translate.AllTranslation[l].PageTop,
		render.H(" ("),
		render.No("a", render.A("href", "/about/"),
			translate.AllTranslation[l].AboutTextLink),
		render.H(")"))
}

// A header to indicated taht this page is currently in development.
func InDevHeader(l language.Language) render.Node {
	return render.N("div.subHeader", translate.AllTranslation[l].InDev)
}
