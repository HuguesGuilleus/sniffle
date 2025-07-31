package eu_eca

import (
	"fmt"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/HuguesGuilleus/sniffle/tool/sch"
)

var annualReportsType = sch.Array(sch.Map(
	sch.FieldSR("Title", sch.NotEmptyString()),
	sch.FieldSR("Description", sch.AnyString()),
	sch.FieldSR("DocSetID", sch.Or(sch.String(""), sch.StrictPositiveStringInt())),
	sch.FieldSR("ImageUrl", sch.Or(sch.Nil(), sch.Regexp(`^(/[\w.\- ,]+)+\.(png|jpg|jpeg)$`))).Comment("Path of image, the host is https://www.eca.europa.eu/"),
	sch.FieldSR("PublicationDate", sch.Time("1/2/2006 15:04:05 PM")),
	sch.FieldSR("ReportLandingPageUrl", sch.Regexp(`^/\w\w/publications/\x{200B}?[\w-]*$`)),
	sch.FieldSR("ReportUrl", sch.And(sch.URL("https://www.eca.europa.eu/**"), sch.Regexp(`\.pdf$|\.PDF$`))),
	sch.FieldSR("Languages", sch.Array(sch.EnumString(
		"BG", "CS", "DA", "DE", "EL", "EN", "ES", "ET", "FI", "FR", "GA", "HR", "HU", "IT", "LT", "LV", "MT", "NL", "PL", "PT", "RO", "SK", "SL", "SV", "RU",
	))).Assert(`xxx`, func(_ map[string]any, field any) error {
		for _, l := range field.([]any) {
			if l == "FR" || l == "EN" {
				return nil
			}
		}
		return fmt.Errorf("Neighter French nor English is not present, langs:%q", field)
	}),
	sch.FieldSR("DocTypes", sch.Nil()),
	sch.FieldSR("IsOpenDataAvailable", sch.False()),
))

var schemaPage = func() []byte {
	l := language.AllEnglish
	title := "European court of auditors reports crawling method"
	description := "Our usage of the https://www.eca.europa.eu/ website to crawl report."
	return render.Merge(render.Na("html", "lang", "en").N(
		render.N("head",
			component.HeadBegin,
			render.N("title", title),
			render.Na("meta", "name", "description").A("content", description),
		),
		render.N("body.edito",
			component.TopHeader(l),
			component.InDevHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerID", "$/eu/eca/report"), ///////
					render.N("div.headerID", "report"),          ///////
					render.N("div.headerID", "schema"),
				),
				render.N("div.headerTitle", title),
			),
			render.N("main.wt.wide",
				component.Toc(l),
				render.N("div.wc",
					render.N("div.summary", description, " It is full empiric so be careful!"),

					render.N("h1", "Report index"),
					render.N("pre.sch",
						annualReportsType.HTML(""),
					),
				),
			),
			component.Footer(l, component.JsSchema|component.JsToc),
		),
	))
}()
