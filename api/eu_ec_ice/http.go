package eu_ec_ice

import (
	"context"
	"fmt"
	"html/template"
	"net/url"
	"slices"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"time"
)

const (
	indexURL  = "https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/0/0"
	detailURL = "https://register.eci.ec.europa.eu/core/api/register/details/%d/%06d"
	logoURL   = "https://register.eci.ec.europa.eu/core/api/register/logo/%d"
)

type ICEOut struct {
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
	// The based language used to write the ICE text.
	DescriptionOriginalLangage language.Language

	// process

	// Members

	TotalSignature uint
	Signature      map[country.Country]uint

	Image []byte
}
type Description struct {
	Title     string
	Website   *url.URL
	Objective template.HTML
	Annex     template.HTML
	Treaty    string
}

func Do(ctx context.Context, t *tool.Tool) {
	iceSlice, err := Fetch(ctx, t)
	if err != nil {
		t.Error("err", "err", err.Error())
		return
	}

	for _, ice := range iceSlice {
		for _, l := range t.Languages {
			renderOne(t, ice, l)
		}
	}
}

func Fetch(ctx context.Context, fetcher tool.Fetcher) ([]*ICEOut, error) {
	items, err := fetchIndex(ctx, fetcher)
	if err != nil {
		return nil, err
	}

	iceSlice := make([]*ICEOut, 0, len(items))
	for _, info := range items {
		ice, err := fetchDetail(ctx, fetcher, info)
		if err != nil {
			return nil, err
		}
		iceSlice = append(iceSlice, ice)
	}

	return iceSlice, nil
}

type indexItem struct {
	year   int
	number int
	logoID int
}

// Get all ICE item to after get all details.
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

func fetchDetail(ctx context.Context, fetcher tool.Fetcher, info indexItem) (*ICEOut, error) {
	ice := &ICEOut{
		Year:        info.year,
		Number:      info.number,
		Description: make(map[language.Language]*Description),
		Signature:   make(map[country.Country]uint),
	}

	dto := &struct {
		Status     string `json:"status"`
		LastUpdate string `json:"latestUpdateDate"`
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
			Entry []struct {
				Country country.Country `json:"countryCodeType"`
				Total   uint            `json:"total"`
			} `json:"entry"`
		} `json:"sosReport"`
	}{}
	if err := tool.FetchGETJSON(ctx, fetcher, fmt.Sprintf(detailURL, info.year, info.number), &dto); err != nil {
		return nil, err
	}

	t, err := time.Parse("02/01/2006 15:04", dto.LastUpdate)
	if err != nil {
		return nil, fmt.Errorf("cannot parse last update %w", err)
	}
	ice.LastUpdate = t

	ice.Status = dto.Status

	categories := make([]string, 0, len(dto.Categories))
	for _, entry := range dto.Categories {
		categories = append(categories, entry.CategoryType)
	}
	slices.Sort(categories)
	ice.Categorie = slices.Compact(categories)

	for _, desc := range dto.Description {
		ice.Description[desc.Language] = &Description{
			Title:     desc.Title,
			Website:   tool.ParseURL(desc.Website),
			Objective: tool.SecureHTML(desc.Objective),
			Annex:     tool.SecureHTML(desc.Annex),
			Treaty:    desc.Treaty,
		}
		if desc.Original {
			ice.DescriptionOriginalLangage = desc.Language
		}
	}

	for _, entry := range dto.Signatures.Entry {
		ice.Signature[entry.Country] = entry.Total
		ice.TotalSignature += entry.Total
	}

	// Image
	if info.logoID != 0 {
		img, err := fetcher.FetchGET(ctx, fmt.Sprintf(logoURL, info.logoID))
		if err != nil {
			return nil, err
		}
		ice.Image = img
	}

	return ice, nil
}

func (ice *ICEOut) GetOriginalDescription() *Description {
	return ice.Description[ice.DescriptionOriginalLangage]
}
