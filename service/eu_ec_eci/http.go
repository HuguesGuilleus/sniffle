package eu_ec_eci

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"slices"
	"sniffle/front/component"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"sniffle/tool/securehtml"
	"time"
)

const (
	indexURL  = "https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/0/0"
	detailURL = "https://register.eci.ec.europa.eu/core/api/register/details/%d/%06d"
	logoURL   = "https://register.eci.ec.europa.eu/core/api/register/logo/%d"
)

type ECIOut struct {
	// Identifier
	Year   int
	Number int

	// Last update of information
	LastUpdate time.Time

	Status string

	// A sorted and uniq slice of categories
	Categorie []string

	// Description text in all language
	Description map[language.Language]*Description
	// The based language used to write the ECI text.
	DescriptionOriginalLangage language.Language

	// process

	// Members

	TotalSignature uint
	Signature      map[country.Country]uint
	// Date of the last organisators signatures update.
	SignaturesUpdate time.Time

	Image []byte
}
type Description struct {
	Title     string
	PlainDesc string
	Website   *url.URL
	Objective template.HTML
	Annex     template.HTML
	Treaty    string
}

func Do(ctx context.Context, t *tool.Tool) {
	eciSlice, err := Fetch(ctx, t)
	if err != nil {
		t.Error("err", "err", err.Error())
		return
	}

	component.RedirectIndex(t, "/eu/ec/eci/")

	for _, eci := range eciSlice {
		component.RedirectIndex(t, fmt.Sprintf("/eu/ec/eci/%d/%d/", eci.Year, eci.Number))
		for _, l := range t.Languages {
			renderOne(t, eci, l)
		}
	}
}

func Fetch(ctx context.Context, fetcher tool.Fetcher) ([]*ECIOut, error) {
	items, err := fetchIndex(ctx, fetcher)
	if err != nil {
		return nil, err
	}

	eciSlice := make([]*ECIOut, 0, len(items))
	for _, info := range items {
		eci, err := fetchDetail(ctx, fetcher, info)
		if err != nil {
			return nil, err
		}
		eciSlice = append(eciSlice, eci)
	}

	return eciSlice, nil
}

type indexItem struct {
	year   int
	number int
	logoID int
}

// Get all ECI item to after get all details.
func fetchIndex(ctx context.Context, fetcher tool.Fetcher) ([]indexItem, error) {
	dto := struct {
		Entries []struct {
			Year   int `json:"year,string"`
			Number int `json:"number,string"`
			Logo   *struct {
				Id int `json:"id"`
			} `json:"logo"`
		} `json:"entries"`
	}{}
	if err := tool.FetchGETJSON(ctx, fetcher, indexURL, &dto); err != nil {
		return nil, err
	}

	items := make([]indexItem, len(dto.Entries))
	for i, dtoEntry := range dto.Entries {
		logoID := 0
		if dtoEntry.Logo != nil {
			logoID = dtoEntry.Logo.Id
		}
		items[i] = indexItem{
			year:   dtoEntry.Year,
			number: dtoEntry.Number,
			logoID: logoID,
		}

	}

	return items, nil
}

func fetchDetail(ctx context.Context, fetcher tool.Fetcher, info indexItem) (*ECIOut, error) {
	eci := &ECIOut{
		Year:        info.year,
		Number:      info.number,
		Description: make(map[language.Language]*Description),
		Signature:   make(map[country.Country]uint),
	}

	dto := &struct {
		Status     string  `json:"status"`
		LastUpdate dtoTime `json:"latestUpdateDate"`
		Categories []struct {
			CategoryType string `json:"categoryType"`
		} `json:"categories"`
		Description []struct {
			Original  bool              `json:"original"`
			Language  language.Language `json:"languageCode"`
			Title     string            `json:"title"`
			Website   string            `json:"website"`
			Objective string            `json:"objectives"`
			Annex     string            `json:"annexText"`
			Treaty    string            `json:"treaties"`
		} `json:"linguisticVersions"`
		Signatures struct {
			UpdateDate dtoDate `json:"updateDate"`
			Entry      []struct {
				Country country.Country `json:"countryCodeType"`
				Total   uint            `json:"total"`
			} `json:"entry"`
		} `json:"sosReport"`
	}{}
	if err := tool.FetchGETJSON(ctx, fetcher, fmt.Sprintf(detailURL, info.year, info.number), &dto); err != nil {
		return nil, err
	}

	eci.LastUpdate = dto.LastUpdate.Time
	eci.Status = dto.Status

	categories := make([]string, 0, len(dto.Categories))
	for _, entry := range dto.Categories {
		categories = append(categories, entry.CategoryType)
	}
	slices.Sort(categories)
	eci.Categorie = slices.Compact(categories)

	for _, desc := range dto.Description {
		eci.Description[desc.Language] = &Description{
			Title:     desc.Title,
			PlainDesc: securehtml.Text(desc.Objective, 200),
			Website:   securehtml.ParseURL(desc.Website),
			Objective: securehtml.Secure(desc.Objective),
			Annex:     securehtml.Secure(desc.Annex),
			Treaty:    desc.Treaty,
		}
		if desc.Original {
			eci.DescriptionOriginalLangage = desc.Language
		}
	}

	eci.SignaturesUpdate = dto.Signatures.UpdateDate.Time
	for _, entry := range dto.Signatures.Entry {
		eci.Signature[entry.Country] = entry.Total
		eci.TotalSignature += entry.Total
	}

	// Image
	if info.logoID != 0 {
		img, err := fetcher.FetchGET(ctx, fmt.Sprintf(logoURL, info.logoID))
		if err != nil {
			return nil, err
		}
		eci.Image = img
	}

	return eci, nil
}

func (eci *ECIOut) GetOriginalDescription() *Description {
	return eci.Description[eci.DescriptionOriginalLangage]
}

type dtoDate struct {
	Time time.Time
}

func (dto *dtoDate) UnmarshalText(data []byte) error {
	t, err := time.ParseInLocation("02/01/2006", string(data), render.DateZone)
	if err != nil {
		return err
	}
	dto.Time = t
	return nil
}

type dtoTime struct {
	Time time.Time
}

func (dto *dtoTime) UnmarshalText(data []byte) error {
	t, err := time.Parse("02/01/2006 15:04", string(data))
	if err != nil {
		return err
	}
	dto.Time = t
	return nil
}
