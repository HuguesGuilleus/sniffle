package eu_eca

import "sniffle/tool/sch"

var annualReportsType = sch.Array(sch.Map(
	sch.FieldSR("Title", sch.NotEmptyString()),
	sch.FieldSR("Description", sch.AnyString()),
	sch.FieldSR("DocSetID", sch.Or(sch.String(""), sch.StrictPositiveStringInt())),
	sch.FieldSR("ImageUrl", sch.Or(sch.Nil(), sch.AnyString())), // path .jpg | .png
	sch.FieldSR("PublicationDate", sch.Time("1/2/2006 15:04:05 PM")),
	sch.FieldSR("ReportLandingPageUrl", sch.AnyString()),
	sch.FieldSR("ReportUrl", sch.URL("https://www.eca.europa.eu/**")), // check .pdf
	sch.FieldSR("Languages", sch.Array(sch.EnumString(
		"BG", "CS", "DA", "DE", "EL", "EN", "ES", "ET", "FI", "FR", "GA", "HR", "HU", "IT", "LT", "LV", "MT", "NL", "PL", "PT", "RO", "SK", "SL", "SV",
	))),
	sch.FieldSR("DocTypes", sch.Nil()),
	sch.FieldSR("IsOpenDataAvailable", sch.False()),
))
