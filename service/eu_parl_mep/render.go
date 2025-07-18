package eu_parl_mep

import (
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/front/translate"
	"sniffle/tool/render"
)

func idNamespace(l language.Language) render.Node {
	tr := translate.T[l]
	return render.N("div.headerID",
		component.HomeAnchor(l),
		render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"), " / ",
		render.Na("a", "href", l.Path("/eu/parl/")).A("title", "!parelement européen").N("parl"), " / ",
		render.Na("a", "href", l.Path("/eu/parl/mep/")).A("title", "!Member of european parlement").N("mep"),
	)
}
