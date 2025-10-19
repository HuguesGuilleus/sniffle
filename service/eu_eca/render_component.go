package eu_eca

import (
	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func headerSupECA(l language.Language) render.Node {
	tr := translate.T[l]
	return render.N("div.headerID",
		component.HomeAnchor(l),
		render.Na("a", "href", l.Path("/eu/")).A("title", tr.EU.Name).N("eu"),
		" / ",
		render.Na("a", "href", l.Path("/eu/eca/")).A("title", tr.EU_ECA.Name).N("eca"),
	)
}

func renderReport(l language.Language, r *report) render.Node {
	ECA := translate.T[l].EU_ECA
	desc := r.Description[l]
	main := render.N("",
		render.N("div.itemTitle.st", desc.Title),
		render.IfS(desc.Description != "",
			render.N("div.itemDesc", desc.Description),
		),
	)

	return render.N("div.si.bigItem",
		r.Image.Render("../report/"+r.ImageHash, ""),
		render.N("div",
			render.N("span.tag", r.PubDate),
			render.N("span.tag", "$$$", r.Type), /////////////////
		),
		render.IfS(r.Langs[l], main),
		render.If(!r.Langs[l], func() render.Node {
			return render.Na("div.tr", "lang", l.String()).N(main)
		}),
		render.If(desc.ReportPage != nil || desc.PDFURL != nil, func() render.Node {
			return render.N("div.boxFlex.boxFlexTop",
				render.If(desc.ReportPage != nil, func() render.Node {
					return render.Na("a.doc", "href", desc.ReportPage.String()).
						N(ECA.ReportPage, " ~>")
				}),
				render.If(desc.PDFURL != nil, func() render.Node {
					return render.Na("a.doc", "href", desc.PDFURL.String()).
						A("title", ECA.ReportPDFHelp).
						N(ECA.ReportPDFHelp, " ~>")
				}),
			)
		}),
	)
}
