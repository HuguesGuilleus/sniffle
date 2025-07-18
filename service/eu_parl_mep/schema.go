package eu_parl_mep

import (
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/tool/render"
	"sniffle/tool/sch"
)

var metaType = sch.Map(
	sch.FieldSR("id", sch.StrictPositiveInt()),
	sch.FieldSR("uri", sch.URL("https://data.europarl.europa.eu/dataset/meps_**")),
	sch.FieldSR("urlDataService", sch.Nil()),
	sch.FieldSR("leg", sch.Or(
		sch.Nil(),
		sch.PositiveInt(),
	)),
	sch.FieldSR("status", sch.String("PUBLISHED")),
	sch.FieldSR("type", sch.String("SERIE")),
	sch.FieldSR("frequency", sch.String("dataset.frequency.irreg")),
	sch.FieldSR("creationDate", sch.StrictPositiveInt()),
	sch.FieldSR("updateDate", sch.StrictPositiveInt()),
	sch.FieldSR("coverDateFrom", sch.AnyInt()),
	sch.FieldSR("coverDateTo", sch.Or(
		sch.Nil(),
		sch.StrictPositiveInt(),
	)),
	sch.FieldSR("views", sch.PositiveInt()),
	sch.FieldSR("download", sch.PositiveInt()),
	sch.FieldSR("locked", sch.False()),
	sch.FieldSR("slug", sch.Nil()),
	sch.FieldSR("reuses", sch.EmptyArray()),
	sch.FieldSR("odpDatasetVersions", sch.Array(sch.MapExtra(
		sch.FieldSR("versionLabel", sch.StrictPositiveStringInt()),
	))),
)

var schemaPage = func() []byte {
	l := language.AllEnglish
	title := "Members of european Parlement crawling method"
	description := "Our methodology to fetch members of european Parlement data"

	return render.Merge(render.Na("html", "lang", "en").N(
		render.N("head",
			component.HeadBegin,
			render.N("title", title),
			render.Na("meta", "name", "description").A("content", description),
		),
		render.N("body.edito",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(l),
					render.N("div.headerID", "schema"),
				),
				render.N("div.headerTitle", title),
			),
			render.N("main.wt.wide",
				component.Toc(l),
				render.N("div.wc",

					render.N("div.summary", "!!!dev..."),

					render.N("h1", "!!!meta data"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", "").N("https://data.europarl.europa.eu/OdpDatasetService/Datasets/members-of-the-european-parliament-meps-parliamentary-term{term}"), "\n\n",
						"200 OK\n",
						"Content-Type: application/json\n\n",
						metaType.HTML("")),
				),
			),
			component.Footer(l, component.JsSchema|component.JsToc),
		),
	))
}()
