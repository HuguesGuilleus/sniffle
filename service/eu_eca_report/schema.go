package eu_eca_report

import (
	"math"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/HuguesGuilleus/sniffle/tool/sch"
)

var indexTypes = sch.Array(sch.Map(
	sch.FieldSO("Description", sch.AnyString()),
	sch.FieldSO("DocSetID", sch.Or(sch.String(""), sch.IntervalStringInt(1, math.MaxInt64))),
	sch.FieldSO("DocTypes", sch.Nil()),
	sch.FieldSO("ImageUrl", sch.Or(
		sch.Nil(),
		sch.Regexp(`^/.+`),
	)),
	sch.FieldSO("IsOpenDataAvailable", sch.False()),
	sch.FieldSO("Languages", sch.Array(sch.EnumString("BG", "CS", "DA", "DE", "EL", "EN", "ES", "ET", "FI", "FR", "GA", "HR", "HU", "IT", "LT", "LV", "MT", "NL", "PL", "PT", "RO", "SK", "SL", "SV"))),
	sch.FieldSO("PublicationDate", sch.Time("1/2/2006 15:04:05 PM")),
	sch.FieldSO("ReportLandingPageUrl", sch.Regexp(`^/../publications/[\w\d_-]+$`)),
	sch.FieldSO("ReportUrl", sch.Regexp(`https://www\.eca\.europa\.eu/.+\.(pdf|PDF)`)),
	sch.FieldSO("Title", sch.NotEmptyString()),
))

var schemaPage = render.Merge(render.Na("html", "lang", "en").N(
	render.N("head", component.HeadBegin,
		render.N("title", "Schema"),
		render.Na("meta", "name", "description").A("content", "Our usage of the API to get the reports."),
	),
	render.N("body.edito",
		component.TopHeader(language.AllEnglish),
		render.N("header",
			render.N("div.headerID",
				component.HomeAnchor(language.AllEnglish),
				render.Na("a", "href", "/eu/").A("title", translate.T[language.AllEnglish].EU.Name).N("eu"), " / ",
				render.Na("a", "href", "/eu/eca/").A("title", translate.T[language.AllEnglish].EU_ECA.Name).N("eca"),
			),
			render.N("div.headerID", "report / schema"),
			render.N("h1", "Schema of ECA"),
		),
		render.N("div.wt.wide",
			component.Toc(language.AllEnglish),
			render.N("main.wc",
				render.N("pre.sch",
					"POST https://www.eca.europa.eu/_vti_bin/ECA.Internet/DocSetService.svc/SearchDocs\n",
					"Accept: application/json\n",
					"Content-Type: application/json\n\n",
					"{\n",
					`	"searchInput":{`, "\n",
					`		"RowLimit": 10000,`, "\n",
					`		"Filter": "ECADocType=\"Annual report\""`, "\n",
					`		"Filter": "ECADocType=\"Specific Annual Report\""`, "\n",
					`		"Filter": "ECADocType=\"Review\""`, "\n",
					`		"Filter": "ECADocType=\"Activity Report\""`, "\n",
					`		"Filter": "ECADocType=\"Journal\""`, "\n",
					`		"Lang": "English"`, "\n",
					`		"LangCode": "EN"`, "\n",
					`		`, "\n",
					`	}`, "\n",
					"}\n",
					"\n",
					"200 OK\n",
					"Content-Type: application/json\n\n",
					indexTypes.HTML(""),
				),
			),
		),
	),
	component.Footer(language.AllEnglish, component.JsSchema|component.JsToc),
))
