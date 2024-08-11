package eu_ec_eci

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"slices"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"sniffle/tool/securehtml"
	"strconv"
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

	ImageName   string
	ImageWidth  string
	ImageHeight string
	ImageData   []byte
}
type Description struct {
	Title       string
	PlainDesc   string
	SupportLink string
	Website     *url.URL
	Objective   template.HTML
	Annex       template.HTML
	Treaty      string
}

type indexItem struct {
	year   int
	number int
	logoID int
}

// Get all ECI item to after get all details.
func fetchIndex(ctx context.Context, t *tool.Tool) []indexItem {
	dto := struct {
		Entries []struct {
			Year   int `json:"year,string"`
			Number int `json:"number,string"`
			Logo   *struct {
				Id int `json:"id"`
			} `json:"logo"`
		} `json:"entries"`
	}{}
	if tool.FetchJSON(ctx, t, indexURL, &dto) {
		return nil
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

	return items
}

func fetchDetail(ctx context.Context, t *tool.Tool, info indexItem) *ECIOut {
	dto := &struct {
		Status     string  `json:"status"`
		LastUpdate dtoTime `json:"latestUpdateDate"`
		Categories []struct {
			CategoryType string `json:"categoryType"`
		} `json:"categories"`
		Description []struct {
			Original    bool              `json:"original"`
			Language    language.Language `json:"languageCode"`
			Title       string            `json:"title"`
			SupportLink string            `json:"supportLink"`
			Website     string            `json:"website"`
			Objective   string            `json:"objectives"`
			Annex       string            `json:"annexText"`
			Treaty      string            `json:"treaties"`
		} `json:"linguisticVersions"`
		Signatures struct {
			UpdateDate dtoDate `json:"updateDate"`
			Entry      []struct {
				Country country.Country `json:"countryCodeType"`
				Total   uint            `json:"total"`
			} `json:"entry"`
		} `json:"sosReport"`
	}{}

	if tool.FetchJSON(ctx, t, fmt.Sprintf(detailURL, info.year, info.number), &dto) {
		return nil
	}

	eci := &ECIOut{
		Year:        info.year,
		Number:      info.number,
		LastUpdate:  dto.LastUpdate.Time,
		Status:      dto.Status,
		Description: make(map[language.Language]*Description),
		Signature:   make(map[country.Country]uint),
	}

	categories := make([]string, 0, len(dto.Categories))
	for _, entry := range dto.Categories {
		categories = append(categories, entry.CategoryType)
	}
	slices.Sort(categories)
	eci.Categorie = slices.Compact(categories)

	for _, desc := range dto.Description {
		eci.Description[desc.Language] = &Description{
			Title:       desc.Title,
			PlainDesc:   securehtml.Text(desc.Objective, 200),
			SupportLink: desc.SupportLink,
			Website:     securehtml.ParseURL(desc.Website),
			Objective:   securehtml.Secure(desc.Objective),
			Annex:       securehtml.Secure(desc.Annex),
			Treaty:      desc.Treaty,
		}
		if desc.Original {
			eci.DescriptionOriginalLangage = desc.Language
		}
	}

	if eci.DescriptionOriginalLangage == language.Invalid {
		t.Warn("noDescription", "year", eci.Year, "nb", eci.Number)
	} else {
		for _, l := range t.Languages {
			if eci.Description[l] == nil {
				eci.Description[l] = eci.GetOriginalDescription()
			}
		}
	}

	eci.SignaturesUpdate = dto.Signatures.UpdateDate.Time
	for _, entry := range dto.Signatures.Entry {
		eci.Signature[entry.Country] = entry.Total
		eci.TotalSignature += entry.Total
	}

	eci.fetchImage(ctx, t, info.logoID)

	return eci
}

func (eci *ECIOut) fetchImage(ctx context.Context, t *tool.Tool, logoID int) {
	if logoID == 0 {
		return
	}

	data := tool.FetchAll(ctx, t, fmt.Sprintf(logoURL, logoID))
	if len(data) == 0 {
		return
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width == 0 || config.Height == 0 {
		return
	}

	switch format {
	case "png":
		eci.ImageName = "logo.png"
	case "jpeg":
		eci.ImageName = "logo.jpg"
	default:
		t.Warn("fetchImage", "err", "unknown format", "format", format)
		return
	}

	eci.ImageWidth = strconv.Itoa(config.Width)
	eci.ImageHeight = strconv.Itoa(config.Height)
	eci.ImageData = data
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
