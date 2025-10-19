package eu_eca

import (
	"fmt"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

func renderIndexByYear(t *tool.Tool, l language.Language, year int, reports []*report) {
	tr := translate.T[l]
	ECA := tr.EU_ECA
	basePath := fmt.Sprintf("/eu/eca/%d/", year)
	title := fmt.Sprintf(ECA.INDEX_BY_YEAR.Title, year)

	t.WriteFile(l.Path(basePath), render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, basePath, title, ECA.INDEX_BY_YEAR.Desc),
		render.N("body",
			component.InDevHeader(l), /////////////////////
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					headerSupECA(l),
					render.N("div.headerID", render.Int(year)),
				),
				render.N("div.headerTitle", title),
				component.HeaderLangs(translate.Langs, l, ""),
			),
			render.N("main.w",
				render.N("div.bigInfo",
					render.N("div.bigInfoMeta", ECA.INDEX_BY_YEAR.Count),
					render.N("div.bigInfoMain.bigInfoData", len(reports)),
				),
				component.SearchBlock(l),
				render.N("div.sg", render.S(reports, "", func(r *report) render.Node {
					return renderReport(l, r)
				})),
			),
			component.Footer(l, component.JsSearch),
		),
	)))
}
