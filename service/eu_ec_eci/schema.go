package eu_ec_eci

import (
	"math"
	"sniffle/front/component"
	"sniffle/tool"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"sniffle/tool/sch"
)

var dateType = sch.Time("02/01/2006")
var timeType = sch.Time("02/01/2006 15:04")

var countriesUpper, countriesUpperDef = sch.Def("CountriesUpper", sch.EnumString("AT", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FR", "GB", "GR", "HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO", "SE", "SI", "SK"))
var countriesLower, countriesLowerDef = sch.Def("CountriesLower", sch.EnumString("at", "be", "bg", "cy", "cz", "de", "dk", "ee", "es", "fi", "fr", "gb", "gr", "hr", "hu", "ie", "it", "lt", "lu", "lv", "mt", "nl", "pl", "pt", "ro", "se", "si", "sk"))
var langs, langsDef = sch.Def("Languages", sch.EnumString("BG", "CS", "DA", "DE", "EL", "EN", "ES", "ET", "FI", "FR", "GA", "HR", "HU", "IT", "LT", "LV", "MT", "NL", "PL", "PT", "RO", "SK", "SL", "SV"))
var statusType, statusTypeDef = sch.Def("Status", sch.EnumString("ANSWERED", "CLOSED", "COLLECTION_START_DATE", "INSUFFICIENT_SUPPORT", "ONGOING", "REGISTERED", "REJECTED", "SUBMITTED", "VERIFICATION", "WITHDRAWN"))

var docType, docTypeDef = sch.Def("Document", sch.Map(
	sch.FieldSR("id", sch.StrictPositiveInt()).Comment(
		"Seel below to fetch document or image.",
	),
	sch.FieldSR("mimeType", sch.MIME(`*/*`)),
	sch.FieldSR("name", sch.NotEmptyString()),
	sch.FieldSR("size", sch.StrictPositiveInt()),
))

var docPDF, docPDFDef = sch.Def("DocumentPDF", sch.And(
	docType,
	sch.MapExtra(sch.FieldSR("mimeType", sch.MIME("application/pdf"))),
))

var docPDFOrMSWord, docPDFOrMSWordDef = sch.Def("DocumentPDFOrMSWord", sch.And(
	docType,
	sch.MapExtra(sch.FieldSR("mimeType", sch.Or(
		sch.MIME("application/pdf"),
		sch.MIME("application/msword"),
		sch.MIME("application/vnd.openxmlformats-officedocument.wordprocessingml.document"),
	))),
))

var docImage, docImageDef = sch.Def("DocumentImage", sch.And(
	docType,
	sch.MapExtra(
		sch.FieldSR("mimeType", sch.Or(sch.MIME("image/png"), sch.MIME("image/jpeg"))),
	),
))

var indexType = sch.Map(
	sch.FieldSR("requests", sch.StrictPositiveInt()),
	sch.FieldSR("registered", sch.StrictPositiveInt()),
	sch.FieldSR("successful", sch.StrictPositiveInt()),
	sch.FieldSR("ongoing", sch.StrictPositiveInt()),
	sch.FieldSR("answered", sch.StrictPositiveInt()),
	sch.FieldSR("all", sch.StrictPositiveInt()),
	sch.FieldSR("recordsFound", sch.StrictPositiveInt()),
	sch.FieldSR("entries", sch.Array(sch.Map(
		sch.FieldSR("id", sch.StrictPositiveInt()),
		sch.FieldSR("year", sch.IntervalStringInt(2012, math.MaxInt64)),
		sch.FieldSR("number", sch.StrictPositiveStringInt()),
		sch.FieldSR("pubRegNum", sch.Regexp(`^ECI\(\d{4}\)\d{6}$`)),
		sch.FieldSR("languageCode", sch.String("EN")),
		sch.FieldSR("lastCall", sch.AnyBool()),
		sch.FieldSR("latestUpdateDate", timeType),
		sch.FieldSO("logo", docImage),
		sch.FieldSR("status", statusType),
		sch.FieldSO("supportLink", sch.AnyURL()),
		sch.FieldSR("title", sch.NotEmptyString()),
		sch.FieldSR("totalSupporters", sch.PositiveInt()),
	))),
)

var eciType = sch.Map(
	sch.FieldSR("id", sch.StrictPositiveInt()),
	sch.FieldSR("comRegNum", sch.Regexp(`^ECI\(\d{4}\)\d{6}$`)),
	sch.FieldSR("status", statusType),
	sch.FieldSR("latestUpdateDate", timeType),
	sch.FieldSR("lastCall", sch.AnyBool()),
	sch.FieldSR("registrationDate", dateType),
	sch.FieldSR("deadline", sch.Or(sch.String(""), dateType)),
	sch.FieldSO("startCollectionDate", dateType),
	sch.FieldSO("earlyClosureDate", sch.Or(sch.String(""), dateType)),
	sch.FieldSR("partiallyRegistered", sch.AnyBool()),
	sch.FieldSR("linguisticVersions", sch.ArrayRange(3, 24, sch.Map(
		sch.FieldSR("original", sch.AnyBool()),
		sch.FieldSR("languageCode", langs),
		sch.FieldSR("title", sch.NotEmptyString()),
		sch.FieldSR("objectives", sch.NotEmptyString()),
		sch.FieldSO("annexText", sch.NotEmptyString()),
		sch.FieldSO("treaties", sch.NotEmptyString()),
		sch.FieldSO("website", sch.AnyURL()),
		sch.FieldSO("supportLink", sch.AnyURL()),
		sch.FieldSR("commissionDecision", sch.Or(
			sch.Map(
				sch.FieldSR("document", docPDF),
			),
			sch.Map(
				sch.FieldSR("celex", sch.NotEmptyString()),
				sch.FieldSO("corrigendum", sch.NotEmptyString()),
				sch.FieldSR("url", sch.URL("http://eur-lex.europa.eu/legal-content/**?uri=*&from=*")),
			),
		)),
		sch.FieldSO("additionalDocument", docPDFOrMSWord),
		sch.FieldSO("draftLegal", docPDFOrMSWord),
	))),
	sch.FieldSO("categories", sch.ArrayMin(1, sch.Map(
		sch.FieldSR("categoryType", sch.EnumString("AGRI", "CULT", "DECO", "DEVCO", "EDU", "EMPL", "ENER", "ENV", "EURO", "JUST", "MARE", "MIGR", "REGIO", "RSH", "SANTE", "SEC", "TRA", "TRADE")),
	))),

	sch.FieldSR("members", sch.ArrayMin(7, sch.And(detailedMember, sch.Map(
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.AnyBool()),
		sch.FieldSR("type", sch.EnumString("MEMBER", "SUBSTITUTE", "REPRESENTATIVE", "OTHER", "DPO", "LEGAL_ENTITY")),
		sch.FieldSO("email", sch.Or(sch.AnyURL(), sch.AnyMail())),
		sch.FieldSO("replacedMember", sch.ArraySize(1, sch.Map(
			sch.FieldSO("email", sch.AnyMail()),
			sch.FieldSO("residenceCountry", countriesLower),
			sch.FieldSR("endDate", dateType),
			sch.FieldSR("fullName", sch.NotEmptyString()),
			sch.FieldSR("privacyApplied", sch.False()),
			sch.FieldSR("startingDate", dateType),
			sch.FieldSR("type", sch.EnumString("MEMBER", "REPRESENTATIVE")),
		))),
		sch.FieldSO("residenceCountry", countriesLower),
		sch.FieldSO("startingDate", dateType),
	)))),

	sch.FieldSR("progress", sch.ArrayMin(2, sch.Map(
		sch.FieldSR("active", sch.AnyBool()),
		sch.FieldSR("name", statusType),
		sch.FieldSO("date", dateType),
		sch.FieldSO("footnoteType", sch.String("COLLECTION_EARLY_CLOSURE")),
	))),
	sch.FieldSR("funding", sch.Or(
		sch.Map(),
		sch.Map(
			sch.FieldSR("lastUpdate", dateType),
			sch.FieldSR("sponsors", sch.ArrayMin(1, sch.Map(
				sch.FieldSR("amount", sch.PositiveFloat()),
				sch.FieldSR("date", dateType),
				sch.FieldSR("name", sch.NotEmptyString()),
				sch.FieldSR("privateSponsor", sch.AnyBool()),
				sch.FieldSR("anonymized", sch.AnyBool()),
				sch.FieldSO("otherSupport", sch.String("Research and Network")).Comment("Found only in ECI 2025/1"),
			))),
			sch.FieldSR("totalAmount", sch.PositiveFloat()),
			sch.FieldSO("document", docPDF),
		),
	)),
	sch.Field(sch.EnumString("submission", "sosReport"), sch.Map(
		sch.FieldSR("totalSignatures", sch.PositiveInt()),
		sch.FieldSO("updateDate", dateType),
		sch.FieldSR("entry", sch.ArrayRange(3, 28, sch.Map(
			sch.FieldSR("countryCodeType", countriesUpper),
			sch.FieldSR("total", sch.PositiveInt()),
			sch.FieldSO("afterSubmission", sch.AnyBool()),
		))),
		sch.FieldSO("footnoteType", sch.String("AFTER_SUBMISSION_CERTIFICATES")),
	), false),
	sch.FieldSO("logo", docImage),
	sch.FieldSO("answer", sch.Map(
		sch.FieldSR("id", sch.StrictPositiveInt()),
		sch.FieldSR("decisionDate", dateType),
		sch.FieldSR("links", sch.ArrayRange(3, 4, sch.Map(
			sch.FieldSO("defaultLanguageCode", sch.String("EN")),
			sch.FieldSR("defaultName", sch.EnumString("ANNEX", "COMMUNICATION", "FOLLOW_UP", "PRESS_RELEASE")),
			sch.FieldSR("defaultLink", sch.URL("http.s://**europa.eu/**")),
			sch.FieldSO("link", sch.ArraySize(24, sch.Map(
				sch.FieldSR("languageCode", langs),
				sch.FieldSR("link", sch.URL("https://**europa.eu/**")),
			))),
		))),
	)),
)

var detailedMember, detailedMemberDef = sch.Def("DetailedMember", sch.Or(
	sch.Map(
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.False()),
		sch.FieldSR("type", sch.String("MEMBER")),
		sch.FieldSO("email", sch.Or(sch.String("email@anonymised"), sch.AnyMail())),
		sch.FieldSO("startingDate", dateType),
		sch.FieldSO("replacedMember", sch.ArraySize(1, sch.Or(
			sch.Map(
				sch.FieldSR("type", sch.String("MEMBER")),
				sch.FieldSR("fullName", sch.NotEmptyString()),
				sch.FieldSR("privacyApplied", sch.False()),
				sch.FieldSR("endDate", dateType),
				sch.FieldSR("startingDate", dateType),
			),
			sch.Map(
				sch.FieldSO("email", sch.AnyMail()),
				sch.FieldSO("residenceCountry", countriesLower),
				sch.FieldSR("endDate", dateType),
				sch.FieldSR("fullName", sch.NotEmptyString()),
				sch.FieldSR("privacyApplied", sch.False()),
				sch.FieldSR("startingDate", dateType),
				sch.FieldSR("type", sch.String("REPRESENTATIVE")),
			),
		))),
	),
	sch.Map(
		sch.FieldSR("type", sch.EnumString("REPRESENTATIVE")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.AnyBool()),
		sch.FieldSO("email", sch.AnyMail()),
		sch.FieldSO("replacedMember", sch.ArraySize(1, sch.Map(
			sch.FieldSO("email", sch.AnyMail()),
			sch.FieldSO("residenceCountry", countriesLower),
			sch.FieldSR("endDate", dateType),
			sch.FieldSR("fullName", sch.NotEmptyString()),
			sch.FieldSR("privacyApplied", sch.False()),
			sch.FieldSR("startingDate", dateType),
			sch.FieldSR("type", sch.EnumString("MEMBER", "REPRESENTATIVE")),
		))),
		sch.FieldSR("residenceCountry", countriesLower),
		sch.FieldSO("startingDate", dateType),
	),
	sch.Map(
		sch.FieldSR("type", sch.EnumString("LEGAL_ENTITY")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.False()),
		sch.FieldSO("email", sch.AnyURL()),
		sch.FieldSR("residenceCountry", countriesLower),
	),
	sch.Map(
		sch.FieldSR("type", sch.EnumString("OTHER")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.False()),
	),
	sch.Map(
		sch.FieldSR("type", sch.EnumString("SUBSTITUTE", "DPO")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.False()),
		sch.FieldSR("email", sch.AnyMail()),
	),
	sch.Map(
		sch.FieldSR("type", sch.EnumString("SUBSTITUTE")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.True()),
	),
))

func renderSchema(t *tool.Tool) {
	lang := language.AllEnglish
	title := "European Citizens' Initiative crawling method"
	description := "Our usage of the https://citizens-initiative.europa.eu/ website to crawl data."

	oneURLTempl := "https://register.eci.ec.europa.eu/core/api/register/details/{year}/{number:%06d}"
	oneURLEx := "https://register.eci.ec.europa.eu/core/api/register/details/2022/000002"

	t.WriteFile("/eu/ec/eci/schema.html", render.Merge(render.Na("html", "lang", "en").N(
		render.N("head",
			component.HeadBegin,
			render.N("title", title),
			render.Na("meta", "name", "description").A("content", description),
		),
		render.N("body.edito",
			component.TopHeader(lang),
			component.InDevHeader(lang),
			render.N("header",
				render.N("div.headerSup",
					idNamespace(lang),
					render.N("div.headerId", "schema"),
				),
				render.N("div.headerTitle", title),
			),
			render.N("main.wt.wide",
				component.Toc,
				render.N("div.wc",
					render.N("div.summary", "Usage of public API https://register.eci.ec.europa.eu/ to get index and details of European Citizens' Initiative. It is full empiric, and official team can do some API change."),

					render.N("h1", "Common types"),
					render.N("h2", "Enumerations types"),
					render.N("pre.sch",
						render.N("span.sch-comment", "// Some types used many times."),
						countriesUpperDef,
						countriesLowerDef,
						langsDef,
						statusTypeDef,
					),
					render.N("h2", "File types"),
					render.N("pre.sch",
						docImageDef,
						docPDFDef,
						docPDFOrMSWordDef,
						docTypeDef,
					),

					render.N("h1", "Get logo or document"),
					render.N("pre.sch",
						render.N("span.sch-comment", "// Fetch a logo:\n"),
						"GET ", render.Na("a.block", "href", "https://register.eci.ec.europa.eu/core/api/register/logo/{logoID}").N("https://register.eci.ec.europa.eu/core/api/register/logo/{logoID}"), "\n",
						render.N("span.sch-comment", "// or to fetch a document:\n"),
						"GET ", render.Na("a.block", "href", "https://register.eci.ec.europa.eu/core/api/register/document/{docID}").N("https://register.eci.ec.europa.eu/core/api/register/document/{docID}"),
						"\n\n200 OK\n",
						`Content-Disposition: attachment; filename ="..."`, "\n",
						"Content-Type: application/octet-stream",
						render.N("hr"),
						render.N("span.sch-comment", "Example: ", render.Na("a", "href", "/eu/ec/eci/2022/2/").N("ECI 2022/2 Fur Free Europe")),
						"\n", render.N("span.sch-comment", "// Financial document:"), "\n",
						"https://register.eci.ec.europa.eu/core/api/register/document/9122",
						"\n", render.N("span.sch-comment", "// Logo:"), "\n",
						"https://register.eci.ec.europa.eu/core/api/register/logo/1979",
					),

					render.N("h1", "Get index range (not used)"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", "https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/{begin}/{end}").N("https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/{begin}/{end}"),

						"\n\n200 OK\n",
						"Content-Type: application/json\n\n",
						"[Same as all index schema below]"),

					render.N("h1", "Get all index"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", indexURL).N(indexURL),
						"\n\n200 OK\n",
						"Content-Type: application/json\n\n",
						indexType.HTML(""),
					),

					render.N("h1", "Get details"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", oneURLTempl).N(oneURLTempl),
						"\n\n# Example: ", render.Na("a.block", "href", oneURLEx).N(oneURLEx),
						"\n\n200 OK\n",
						"Content-Type: application/json\n\n",
						eciType.HTML(""),
					),

					render.N("h2", "Detailed member"),
					render.N("pre.sch", detailedMemberDef),

					render.N("h1", "Thresholds data"),
					render.N("p.noindent",
						"We manualy extract data from ",
						render.Na("a", "href", "https://citizens-initiative.europa.eu/thresholds_en").N("https://citizens-initiative.europa.eu/thresholds_en"),
						". Last check: ", threshold_lastCheck, ".",
					),
				),
			),
			component.Footer(lang, component.JsSchema|component.JsToc),
		),
	)))
}
