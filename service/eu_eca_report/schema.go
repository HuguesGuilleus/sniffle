package eu_eca_report

import (
	"math"

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
