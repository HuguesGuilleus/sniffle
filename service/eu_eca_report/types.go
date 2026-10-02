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
	ID string
	// The type of report.
	// Values is equal to API `ECADocType`.
	Kind            string
	PublicationDate time.Time
	Image           *rimage.Image `json:"-"`
	ImageURL        string
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

func (all *allReport) All() (a []report) {
	a = make([]report, 0,
		len(all.AnnualReport)+
			len(all.SpecificAnnualReport)+
			len(all.Review)+
			len(all.ActivityJournal)+
			len(all.Journal),
	)
	a = append(a, all.AnnualReport...)
	a = append(a, all.SpecificAnnualReport...)
	a = append(a, all.Review...)
	a = append(a, all.ActivityJournal...)
	a = append(a, all.Journal...)
	slices.SortFunc(a, func(a, b report) int {
		return b.PublicationDate.Compare(a.PublicationDate)
	})
	return
}

func (all *allReport) ByYear() iter.Seq2[int, []*report] {
	byYears := make(map[int][]*report)
	for _, r := range all.All() {
		year := r.PublicationDate.Year()
		byYears[year] = append(byYears[year], &r)
	}
	return func(yield func(int, []*report) bool) {
		for year, reports := range byYears {
			slices.SortFunc(reports, func(a, b *report) int {
				return b.PublicationDate.Compare(a.PublicationDate)
			})
			if !yield(year, reports) {
				return
			}
		}
	}
}
