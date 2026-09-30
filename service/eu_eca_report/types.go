package eu_eca_report

import (
	"iter"
	"net/url"
	"slices"
	"time"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/common/rimage"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

type allReport struct {
	AnnualReport         []report
	SpecificAnnualReport []report
	Review               []report
	ActivityJournal      []report
	Journal              []report
}
type report struct {
	ID              string
	PublicationDate time.Time
	Image           *rimage.Image
	L               [language.Len]*reportL
}
type reportL struct {
	L                    language.Language
	Title                string
	Description          render.H
	ReportLandingPageUrl *url.URL
	ReportUrl            *url.URL
}

func (all *allReport) Ids() (allIds []string) {
	allIds = make([]string, 0, 0+
		len(all.AnnualReport)+
		len(all.SpecificAnnualReport)+
		len(all.Review)+
		len(all.ActivityJournal)+
		len(all.Journal),
	)
	for _, r := range all.AnnualReport {
		allIds = append(allIds, r.ID)
	}
	for _, r := range all.SpecificAnnualReport {
		allIds = append(allIds, r.ID)
	}
	for _, r := range all.Review {
		allIds = append(allIds, r.ID)
	}
	for _, r := range all.ActivityJournal {
		allIds = append(allIds, r.ID)
	}
	for _, r := range all.Journal {
		allIds = append(allIds, r.ID)
	}
	return
}

func (all *allReport) All() iter.Seq[report] {
	return func(yield func(report) bool) {
		slices.Values(all.AnnualReport)(yield)
		slices.Values(all.SpecificAnnualReport)(yield)
		slices.Values(all.Review)(yield)
		slices.Values(all.ActivityJournal)(yield)
		slices.Values(all.Journal)(yield)
	}
}
