package eu_ec_ice

import (
	"context"
	"fmt"
	"sniffle/tool"
	"sniffle/tool/country"
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

	// process

	// about text ...

	// Members

	TotalSignature uint
	Signature      map[country.Country]uint

	Image []byte
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
	// ice := new(ICEOut)
	ice := &ICEOut{
		Year:      info.year,
		Number:    info.number,
		Signature: make(map[country.Country]uint),
	}

	dto := &struct {
		Status     string `json:"status"`
		LastUpdate string `json:"latestUpdateDate"`
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

	ice.Status = dto.Status

	t, err := time.Parse("02/01/2006 15:04", dto.LastUpdate)
	if err != nil {
		return nil, fmt.Errorf("cannot parse last update %w", err)
	}
	ice.LastUpdate = t

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
