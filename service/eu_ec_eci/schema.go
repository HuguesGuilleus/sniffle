package eu_ec_eci

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/HuguesGuilleus/sniffle/tool/sch"
)

var dateType = sch.Time("02/01/2006")
var timeType = sch.Time("02/01/2006 15:04")

var countriesUpper, countriesUpperDef = sch.Def("CountriesUpper", sch.EnumString("AT", "BE", "BG", "CY", "CZ", "DE", "DK", "EE", "ES", "FI", "FR", "GB", "GR", "HR", "HU", "IE", "IT", "LT", "LU", "LV", "MT", "NL", "PL", "PT", "RO", "SE", "SI", "SK"))
var countriesLower, countriesLowerDef = sch.Def("CountriesLower", sch.EnumString("at", "be", "bg", "cy", "cz", "de", "dk", "ee", "es", "fi", "fr", "gb", "gr", "hr", "hu", "ie", "it", "lt", "lu", "lv", "mt", "nl", "pl", "pt", "ro", "se", "si", "sk"))
var langs, langsDef = sch.Def("Languages", sch.EnumString("BG", "CS", "DA", "DE", "EL", "EN", "ES", "ET", "FI", "FR", "GA", "HR", "HU", "IT", "LT", "LV", "MT", "NL", "PL", "PT", "RO", "SK", "SL", "SV"))
var statusType, statusTypeDef = sch.Def("Status", sch.EnumString("ANSWERED", "CLOSED", "COLLECTION_START_DATE", "INSUFFICIENT_SUPPORT", "ONGOING", "REGISTERED", "REJECTED", "SUBMITTED", "VERIFICATION", "WITHDRAWN"))

var docType, docTypeDef = sch.Def("Document", sch.Map(
	sch.FieldSR("id", sch.StrictPositiveInt()).Comment(
		"See below to fetch document or image.",
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

var docApplication, docApplicationDef = sch.Def("DocumentPDFOrMSWord", sch.And(
	docType,
	sch.MapExtra(sch.FieldSR("mimeType", sch.Or(
		sch.MIME("application/pdf"),
		sch.MIME("application/msword"),
		sch.MIME("application/vnd.openxmlformats-officedocument.wordprocessingml.document"),
		sch.MIME("application/force-download"),
		sch.MIME("application/octet-stream"),
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
		sch.FieldSR("pubRegNum", sch.Regexp(`^ECI\(\d{4}\)\d{6}$`)).Assert(`== "ECI(year)number"`, func(eci map[string]any, field any) error {
			expected := fmt.Sprintf("ECI(%s)%s", eci["year"], eci["number"])
			if field.(string) != expected {
				return fmt.Errorf("%q != %q", field, expected)
			}
			return nil
		}),
		sch.FieldSR("languageCode", sch.String("EN")),
		sch.FieldSR("lastCall", sch.AnyBool()),
		sch.FieldSR("latestUpdateDate", timeType),
		sch.FieldSO("logo", docImage),
		sch.FieldSR("status", statusType),
		sch.FieldSO("supportLink", sch.AnyURL()),
		sch.FieldSR("title", sch.NotEmptyString()),
		sch.FieldSR("totalSupporters", sch.PositiveInt()),
	))).Assert(sch.AssertKey("id", func(eci any) int64 {
		id, _ := eci.(map[string]any)["id"].(json.Number).Int64()
		return id
	})).Assert(sch.AssertKey("pubRegNum", func(eci any) string {
		return eci.(map[string]any)["pubRegNum"].(string)
	})),
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
	sch.FieldSO("logo", docImage),
	sch.BlankField(),
	sch.FieldSR("linguisticVersions", sch.ArrayRange(3, 24, sch.Map(
		sch.FieldSR("original", sch.AnyBool()),
		sch.FieldSR("languageCode", langs),
		sch.FieldSR("title", sch.NotEmptyString()),
		sch.FieldSR("objectives", sch.NotEmptyString()),
		sch.FieldSO("annexText", sch.NotEmptyString()),
		sch.FieldSO("treaties", sch.NotEmptyString()),
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
		sch.FieldSO("website", sch.AnyURL()),
		sch.FieldSO("supportLink", sch.AnyURL()),
	))).Assert(sch.AssertKey(`languageCode`, func(desc any) string {
		return desc.(map[string]any)["languageCode"].(string)
	})).Assert(sch.AssertOnlyOneTrue("original", func(desc any) bool {
		return desc.(map[string]any)["original"].(bool)
	})),
	sch.BlankField(),
	sch.FieldSO("categories", sch.ArrayMin(1, sch.Map(
		sch.FieldSR("categoryType", sch.EnumString("AGRI", "CULT", "DECO", "DEVCO", "EDU", "EMPL", "ENER", "ENV", "EURO", "JUST", "MARE", "MIGR", "REGIO", "RSH", "SANTE", "SEC", "TRA", "TRADE")),
	))).Assert(sch.AssertKey("categoryType", func(category any) string {
		return category.(map[string]any)["categoryType"].(string)
	})),
	sch.BlankField(),

	sch.FieldSR("members", sch.ArrayMin(7, sch.And(detailedMember, sch.Map(
		sch.FieldSR("type", sch.EnumString("MEMBER", "SUBSTITUTE", "REPRESENTATIVE", "OTHER", "DPO", "LEGAL_ENTITY")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSO("residenceCountry", countriesLower),
		sch.FieldSO("startingDate", dateType),
		sch.FieldSR("privacyApplied", sch.AnyBool()),
		sch.FieldSO("email", sch.Or(
			sch.EnumString("anonymised@anonymised", "email@anonymised"),
			sch.AnyMail(),
			sch.AnyURL(),
		)),
		sch.FieldSO("replacedMember", sch.ArraySize(1, sch.Map(
			sch.FieldSR("type", sch.EnumString("MEMBER", "REPRESENTATIVE", "SUBSTITUTE")),
			sch.FieldSR("fullName", sch.NotEmptyString()),
			sch.FieldSR("endDate", dateType),
			sch.FieldSR("startingDate", dateType),
			sch.FieldSR("privacyApplied", sch.False()),
			sch.FieldSO("email", sch.AnyMail()),
			sch.FieldSO("residenceCountry", countriesLower),
		))),
	)))),

	sch.BlankField(),
	sch.FieldSR("progress", sch.ArrayMin(1, sch.Map(
		sch.FieldSR("active", sch.AnyBool()),
		sch.FieldSR("name", statusType),
		sch.FieldSO("date", dateType),
		sch.Assert(`date == null => name == "COLLECTION_START_DATE" && active == false`, func(this map[string]any, _ any) error {
			status := this["name"].(string)
			if this["date"] == nil && status != "COLLECTION_START_DATE" {
				return fmt.Errorf("wrong status when no date: %q", status)
			}
			if this["date"] == nil && this["active"].(bool) {
				return fmt.Errorf("must not be active")
			}
			return nil
		}),
		sch.FieldSO("footnoteType", sch.String("COLLECTION_EARLY_CLOSURE")).Assert(`this.name == "CLOSED"`, func(this map[string]any, field any) error {
			if this["name"] != "CLOSED" {
				return fmt.Errorf("expect this.name == CLOSED, but get %q", this["name"])
			}
			return nil
		}),
	))).Assert(sch.AssertKey("name", func(progress any) string {
		return progress.(map[string]any)["name"].(string)
	})).Assert(sch.AssertOnlyOneTrue("active", func(progress any) bool {
		return progress.(map[string]any)["active"].(bool)
	})),
	sch.Assert(`eci.status == this.progress[with .active].name OR (eci.status == "REGISTERED" AND this.progress[with .active].name == "COLLECTION_START_DATE")`, func(eci map[string]any, _ any) error {
		eciStatus := eci["status"].(string)
		if eciStatus == "REGISTERED" {
			eciStatus = "COLLECTION_START_DATE"
		}
		for _, p := range eci["progress"].([]any) {
			if progress := p.(map[string]any); progress["active"].(bool) {
				if active := progress["name"].(string); active != eciStatus {
					return fmt.Errorf("Not same status: eci.status:%q, active:%q", eciStatus, active)
				}
				break
			}
		}
		return nil
	}),
	sch.BlankField(),
	sch.FieldSR("funding", sch.Or(
		sch.Map(),
		sch.Map(
			sch.FieldSR("lastUpdate", dateType),
			sch.FieldSR("sponsors", sch.ArrayMin(1, sch.Map(
				sch.FieldSR("amount", sch.PositiveFloat()),
				sch.FieldSR("date", dateType),
				sch.FieldSR("name", sch.NotEmptyString()),
				sch.FieldSR("privateSponsor", sch.AnyBool()),
				sch.FieldSR("anonymized", sch.AnyBool()).Assert(`if anonymized => privateSponsor == true && name == "[ANONYMIZED]"`, func(this map[string]any, field any) error {
					if field.(bool) {
						if !this["privateSponsor"].(bool) {
							return fmt.Errorf("The anonymized sponsor is not private")
						}
						if name := this["name"].(string); name != "[ANONYMIZED]" {
							return fmt.Errorf("The sponsor name is not '[ANONYMIZED]', is %q", name)
						}
					}
					return nil
				}),
				sch.FieldSO("otherSupport", sch.EnumString(
					"Film Screening Exhibition (Organisation, Design and Set-Up, Communication)",
					"in-kind donation",
					"National Press Conference (Organisation, Speakers, Communication)",
					"Research and Network",
					"Traveling Exhibition (Organisation, Design and Set-Up, Conference, Communication)",
				)).Comment("Found only in ECI 2025/1"),
			))),
			sch.FieldSR("totalAmount", sch.PositiveFloat()).Assert(`totalAmount == sum(sponsors[*].amount)`, func(this map[string]any, field any) error {
				totalAmount, _ := field.(json.Number).Float64()
				sum := float64(0)
				for _, sponsors := range this["sponsors"].([]any) {
					f, _ := sponsors.(map[string]any)["amount"].(json.Number).Float64()
					sum += f
				}
				sum = math.Round(sum*100) / 100
				if sum != totalAmount {
					return fmt.Errorf("totalAmount %f != sum(sponsors[*].amount) %f", totalAmount, sum)
				}
				return nil
			}),
			sch.FieldSO("document", docPDF),
		),
	)),
	sch.BlankField(),
	sch.Field(sch.EnumString("sosReport", "submission"), sch.Map(
		sch.FieldSO("updateDate", dateType),
		sch.FieldSR("entry", sch.ArrayRange(3, 28, sch.Map(
			sch.FieldSR("countryCodeType", countriesUpper),
			sch.FieldSR("total", sch.PositiveInt()),
			sch.FieldSO("afterSubmission", sch.AnyBool()),
		))),
		sch.FieldSO("footnoteType", sch.String("AFTER_SUBMISSION_CERTIFICATES")),
		sch.Assert(`footnoteType exist <=> some(entry[$].afterSubmission == true)`, func(this map[string]any, _ any) error {
			someAfterSubmission := false
			for _, entry := range this["entry"].([]any) {
				if after := entry.(map[string]any)["afterSubmission"]; after != nil {
					someAfterSubmission = someAfterSubmission || after.(bool)
				}
			}
			footnoteType := (this["footnoteType"] != nil)
			if someAfterSubmission != footnoteType {
				return fmt.Errorf("someAfterSubmission %t != footnoteType existance %t", someAfterSubmission, footnoteType)
			}
			return nil
		}),
		sch.FieldSR("totalSignatures", sch.PositiveInt()).Assert(`totalSignatures == sum(entry[*].total without .afterSubmission==true)`, func(this map[string]any, field any) error {
			totalSignatures, _ := field.(json.Number).Int64()
			sum := int64(0)
			for _, entry := range this["entry"].([]any) {
				if after := entry.(map[string]any)["afterSubmission"]; after != nil && after.(bool) == true {
					continue
				}
				total, _ := entry.(map[string]any)["total"].(json.Number).Int64()
				sum += total
			}
			if totalSignatures != sum {
				return fmt.Errorf("totalSignatures %d != sum(entry[*].total) %d", totalSignatures, sum)
			}
			return nil
		}),
		sch.FieldSO("onlineSosUpdateDate", timeType),
	), false).Comment("sosReport when collect is not validated, else submission field."),
	sch.BlankField(),
	sch.FieldSO("preAnswer", sch.Map(
		sch.FieldSR("links", sch.ArraySize(1,
			sch.Map(
				sch.FieldSR("defaultLanguageCode", sch.String("EN")),
				sch.FieldSR("defaultName", sch.String("EXAMIN_STEPS")),
				sch.FieldSR("defaultLink", sch.URL("https://citizens-initiative.europa.eu/**?")),
			),
		)),
	)).Comment("Only view for ECI 2019/7. View 2025-03-11"),
	sch.Assert(`eci.status==SUBMITED <=> eci.preAnswer={...}`, func(eci map[string]any, _ any) error {
		if (eci["status"] == "SUBMITTED") != (eci["preAnswer"] != nil) {
			return fmt.Errorf("eci.status(actual:%q)==SUBMITED <=> eci.preAnswer(actualNotNull:%t)={...}", eci["status"], eci["preAnswer"] != nil)
		}
		return nil
	}),
	sch.BlankField(),
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
			))).Assert(sch.AssertKey("languageCode", func(link any) string {
				return link.(map[string]any)["languageCode"].(string)
			})),
		))).Assert(sch.AssertKey("defaultName", func(links any) string {
			return links.(map[string]any)["defaultName"].(string)
		})),
	)).Assert(`eci.status == "ANSWERED"`, func(eci map[string]any, _ any) error {
		if eci["status"] != "ANSWERED" {
			return fmt.Errorf("expect SUBMITTED status, but get %q", eci["status"])
		}
		return nil
	}),
)

var detailedMember, detailedMemberDef = sch.Def("DetailedMember", sch.Or(
	sch.Map(
		sch.FieldSR("type", sch.EnumString("SUBSTITUTE", "REPRESENTATIVE", "MEMBER")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.AnyBool()),
		sch.FieldSO("email", sch.Or(
			sch.String("anonymised@anonymised"),
			sch.String("email@anonymised"),
			sch.AnyMail(),
		)),
		sch.FieldSO("startingDate", dateType),
		sch.FieldSO("residenceCountry", countriesLower),
		sch.FieldSO("replacedMember", sch.ArraySize(1, sch.Map(
			sch.FieldSR("type", sch.EnumString("SUBSTITUTE", "REPRESENTATIVE", "MEMBER")),
			sch.FieldSR("fullName", sch.NotEmptyString()),
			sch.FieldSR("privacyApplied", sch.False()),
			sch.FieldSR("startingDate", dateType),
			sch.FieldSR("endDate", dateType),
			sch.FieldSO("email", sch.AnyMail()),
			sch.FieldSO("residenceCountry", countriesLower),
		))),
	),

	sch.Map(
		sch.FieldSR("type", sch.String("DPO")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.False()),
		sch.FieldSR("email", sch.AnyMail()),
	),

	sch.Map(
		sch.FieldSR("type", sch.String("LEGAL_ENTITY")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.False()),
		sch.FieldSO("email", sch.AnyURL()),
		sch.FieldSR("residenceCountry", countriesLower),
	),

	sch.Map(
		sch.FieldSR("type", sch.String("OTHER")),
		sch.FieldSR("fullName", sch.NotEmptyString()),
		sch.FieldSR("privacyApplied", sch.False()),
	),
))

var refusedIndexType = sch.Map(
	sch.FieldSR("requests", sch.PositiveInt()),
	sch.FieldSR("registered", sch.PositiveInt()),
	sch.FieldSR("successful", sch.PositiveInt()),
	sch.FieldSR("ongoing", sch.PositiveInt()),
	sch.FieldSR("answered", sch.PositiveInt()),
	sch.FieldSR("all", sch.PositiveInt()),
	sch.FieldSR("recordsFound", sch.PositiveInt()),
	sch.FieldSR("entries", sch.Array(sch.Map(
		sch.FieldSR("id", sch.PositiveInt()),
		sch.FieldSR("status", sch.String("REJECTED")),
		sch.FieldSR("languageCode", langs),
		sch.FieldSR("title", sch.NotEmptyString()),
		sch.FieldSR("totalSupporters", sch.ConstInt(0)),
		sch.FieldSR("latestUpdateDate", timeType),
		sch.FieldSR("lastCall", sch.False()),
	))),
)

var refusedOneType = sch.Map(
	sch.FieldSR("id", sch.PositiveInt()),
	sch.FieldSR("latestUpdateDate", timeType),
	sch.FieldSR("lastCall", sch.False()),
	sch.FieldSR("refusalDate", dateType),
	sch.FieldSR("refusalDocument", docPDF),
	sch.FieldSR("linguisticVersions", sch.ArraySize(1, sch.Map(
		sch.FieldSR("original", sch.True()),
		sch.FieldSR("languageCode", langs),
		sch.FieldSR("title", sch.NotEmptyString()),
		sch.FieldSR("objectives", sch.NotEmptyString()),
		sch.FieldSO("annexText", sch.NotEmptyString()),
		sch.FieldSR("treaties", sch.NotEmptyString()),
		sch.FieldSO("website", sch.AnyURL()),
		sch.FieldSO("supportLink", sch.AnyURL()),
		sch.Assert(`this.website == tis.supportLink`, func(this map[string]any, _ any) error {
			if this["website"] != this["supportLink"] {
				return fmt.Errorf("website:%q != supportLink:%q", this["website"], this["supportLink"])
			}
			return nil
		}),
		sch.FieldSO("additionalDocument", docApplication).Comment("additionalDocument and draftLegal fields can be equal execpt id."),
		sch.FieldSO("draftLegal", docPDFOrMSWord),
		sch.FieldSR("commissionDecision", sch.Map(
			sch.FieldSR("document", docPDF),
			sch.FieldSO("celex", sch.NotEmptyString()),
		)),
	))),
	sch.Assert(`eci.refusalDocument == eci.linguisticVersions[0].commissionDecision.document`, func(this map[string]any, _ any) error {
		commissionDoc := this["linguisticVersions"].([]any)[0].(map[string]any)["commissionDecision"].(map[string]any)["document"]
		refusalDocument := this["refusalDocument"]
		if !reflect.DeepEqual(refusalDocument, commissionDoc) {
			return fmt.Errorf("refusalDocument:%+v, commissionDoc:%+v", refusalDocument, commissionDoc)
		}
		return nil
	}),
	sch.FieldSO("refusalReasons", sch.ArraySize(1, sch.String("reason.action.registration.reject.competences"))),
)

var schemaPage = func() []byte {
	l := language.AllEnglish
	title := "European Citizens' Initiative crawling method"
	description := "Our usage of the https://citizens-initiative.europa.eu/ website to crawl data."

	oneURLTempl := "https://register.eci.ec.europa.eu/core/api/register/details/{year}/{number:%06d}"
	oneURLTemplId := "https://register.eci.ec.europa.eu/core/api/register/details/{id}"
	oneURLExyear := "https://register.eci.ec.europa.eu/core/api/register/details/2022/000002"
	oneURLExId := "https://register.eci.ec.europa.eu/core/api/register/details/1979"

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
					render.N("div.summary", "Usage of public API https://register.eci.ec.europa.eu/ to get index and details of European Citizens' Initiative. It is full empiric, be careful!"),

					render.N("h1", "Enumerations types"),
					render.N("pre.sch",
						countriesUpperDef, "\n\n",
						countriesLowerDef, "\n\n",
						langsDef, "\n\n",
						statusTypeDef,
					),

					render.N("h1", "Logo or document"),
					render.N("h2", "Types"),
					render.N("pre.sch",
						docTypeDef, "\n\n",
						docImageDef, "\n\n",
						docPDFDef, "\n\n",
						docPDFOrMSWordDef,
					),
					render.N("h2", "Fetch the file"),
					render.N("pre.sch",
						component.SchComment("Fetch a logo:"),
						"GET ", render.Na("a.block", "href", "https://register.eci.ec.europa.eu/core/api/register/logo/{logoID}").N("https://register.eci.ec.europa.eu/core/api/register/logo/{logoID}"), "\n",
						component.SchComment("Fetch a document:"),
						"GET ", render.Na("a.block", "href", "https://register.eci.ec.europa.eu/core/api/register/document/{docID}").N("https://register.eci.ec.europa.eu/core/api/register/document/{docID}"), "\n",
						component.SchComment("Response:"),
						"200 OK\n",
						`Content-Disposition: attachment; filename ="..."`, "\n",
						"Content-Type: application/octet-stream",
						render.N("hr"),
						component.SchComment("Example: ", render.Na("a", "href", "/eu/ec/eci/2022/2/").N("ECI 2022/2 Fur Free Europe")),
						component.SchComment("- Financial document:"),
						"https://register.eci.ec.europa.eu/core/api/register/document/9122\n",
						component.SchComment("- Logo:"),
						"https://register.eci.ec.europa.eu/core/api/register/logo/1979",
					),

					render.N("h1", "Get index"),
					render.N("h2", "Get only a range (not used)"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", "https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/{begin}/{end}").N("https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/{begin}/{end}"),

						"\n\n200 OK\n",
						"Content-Type: application/json\n\n",
						"[Same as all index schema below]"),

					render.N("h2", "Get full index"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", acceptedIndexURL).N(acceptedIndexURL),
						"\n\n200 OK\n",
						"Content-Type: application/json\n\n",
						indexType.HTML(""),
					),

					render.N("h2", "Get refused index"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", refusedIndexURL).N(refusedIndexURL),
						"\n\n200 OK\n",
						"Content-Type: application/json\n\n",
						refusedIndexType.HTML(""),
					),

					render.N("h1", "Get details"),
					render.N("h2", "Fetch"),
					render.N("pre.sch",
						"GET ", render.Na("a.block", "href", oneURLTemplId).N(oneURLTemplId), " ",
						component.SchComment("For all ECI"),
						"GET ", render.Na("a.block", "href", oneURLTempl).N(oneURLTempl), " ",
						component.SchComment("Only for non refused ICE"),
						"\n200 OK\n",
						"Content-Type: application/json\n\n{...}",
						render.N("hr"),
						"# Example: Fur Free Europe\n",
						render.Na("a.block", "href", oneURLExyear).N(oneURLExyear), "\n",
						render.Na("a.block", "href", oneURLExId).N(oneURLExId),
					),
					render.N("h2", "Accepted ECI"),
					render.N("pre.sch", eciType.HTML("")),

					render.N("h2", "Detailed member"),
					render.N("pre.sch", detailedMemberDef),

					render.N("h2", "Refused ECI"),
					render.N("pre.sch", refusedOneType.HTML("")),

					render.N("h1", "Thresholds data"),
					render.N("p.noindent",
						"We manualy extract data from ",
						render.Na("a", "href", "https://citizens-initiative.europa.eu/thresholds_en").N("https://citizens-initiative.europa.eu/thresholds_en"),
						" and legal text. ",
						render.Na("a", "href", "data-threshold/").N("See used data."),
						" Last check: ", threshold_lastCheck, ".",
					),
				),
			),
			component.Footer(l, component.JsSchema|component.JsToc),
		),
	))
}()
