package eu_ec_eci

import (
	"bytes"
	"cmp"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mime"
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

	Timeline []Timeline

	TotalSignature     uint
	ValidatedSignature bool
	Signature          map[country.Country]uint
	// Date of the last organisators paper signatures update.
	// Can be zero
	PaperSignaturesUpdate time.Time
	Threshold             Threshold
	ThresholdPassed       uint

	ImageName   string
	ImageWidth  string
	ImageHeight string
	ImageData   []byte
}
type Description struct {
	Title       string
	PlainDesc   string
	SupportLink *url.URL
	Website     *url.URL
	Objective   render.H
	Annex       render.H
	Treaty      render.H
}
type Timeline struct {
	Date       time.Time
	Status     string
	EarlyClose bool
	Register   *[language.Len]Document
}
type Threshold = map[country.Country]uint
type Document struct {
	URL      *url.URL
	Name     string
	Size     int
	MimeType string
}

type indexItem struct {
	year   int
	number int
	logoID int
}

// Get all ECI item to after get all details.
func fetchIndex(ctx context.Context, t *tool.Tool) (items []indexItem) {
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
	if t.Dev() {
		t.WriteFile("/eu/ec/eci/src.json", tool.FetchAll(ctx, t, indexURL))
	}

	items = make([]indexItem, len(dto.Entries))
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

	return
}

func fetchDetail(ctx context.Context, t *tool.Tool, info indexItem) *ECIOut {
	type signatureDTO struct {
		Country country.Country `json:"countryCodeType"`
		Total   uint
	}
	type docDTO struct {
		Id       int
		Name     string
		Size     int
		MimeType mimeTypeDTO
	}
	dto := &struct {
		Status      string
		LastUpdate  dtoTime `json:"latestUpdateDate"`
		Categories  []struct{ CategoryType string }
		Description []struct {
			Original    bool
			Language    language.Language `json:"languageCode"`
			Title       string
			SupportLink string
			Website     string
			Objective   string `json:"objectives"`
			Annex       string `json:"annexText"`
			Treaty      string `json:"treaties"`
			Register    struct {
				Url      string
				Document *docDTO
			} `json:"commissionDecision"`
		} `json:"linguisticVersions"`
		Progress []struct {
			Status string  `json:"Name"`
			Date   dtoDate `json:"date"`
			Note   string  `json:"footnoteType"`
		} `json:"progress"`
		Signatures struct {
			UpdateDate dtoDate `json:"updateDate"`
			Entry      []signatureDTO
		} `json:"sosReport"`
		Submission struct {
			Entry []signatureDTO
		}
	}{}
	fetchURL := fmt.Sprintf(detailURL, info.year, info.number)
	if t.Dev() {
		t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/src.json", info.year, info.number),
			tool.FetchAll(ctx, t, fetchURL))
	}
	if tool.FetchJSON(ctx, t, fetchURL, &dto) {
		return nil
	}

	eci := &ECIOut{
		Year:        info.year,
		Number:      info.number,
		LastUpdate:  dto.LastUpdate.Time,
		Categorie:   make([]string, 0, len(dto.Categories)),
		Status:      dto.Status,
		Description: make(map[language.Language]*Description),
		Signature:   make(map[country.Country]uint),
	}

	// Categorie
	for _, entry := range dto.Categories {
		eci.Categorie = append(eci.Categorie, entry.CategoryType)
	}
	slices.Sort(eci.Categorie)
	eci.Categorie = slices.Compact(eci.Categorie)

	// Description
	registrationDoc := new([language.Len]Document)
	for _, desc := range dto.Description {
		if desc.SupportLink == desc.Website {
			desc.SupportLink = ""
		}
		eci.Description[desc.Language] = &Description{
			Title:       desc.Title,
			PlainDesc:   securehtml.Text(desc.Objective, 200),
			SupportLink: securehtml.ParseURL(desc.SupportLink),
			Website:     securehtml.ParseURL(desc.Website),
			Objective:   securehtml.Secure(desc.Objective),
			Annex:       securehtml.Secure(desc.Annex),
			Treaty:      securehtml.TextWithURL(desc.Treaty),
		}
		if d := eci.Description[desc.Language]; d.SupportLink != nil && d.SupportLink.Host == "ec.europa.eu" {
			d.SupportLink = nil
		}
		if desc.Original {
			eci.DescriptionOriginalLangage = desc.Language
		}
		if u := securehtml.ParseURL(desc.Register.Url); u != nil {
			registrationDoc[desc.Language] = Document{URL: u}
		} else if desc.Register.Document != nil {
			registrationDoc[desc.Language] = Document{
				URL: &url.URL{
					Scheme: "https",
					Host:   "register.eci.ec.europa.eu",
					Path:   "/core/api/register/document/" + strconv.Itoa(desc.Register.Document.Id),
				},
				Name:     desc.Register.Document.Name,
				MimeType: string(desc.Register.Document.MimeType),
				Size:     desc.Register.Document.Size,
			}
		} else {
			t.Warn("noRegistrationDoc", "year", eci.Year, "nb", eci.Number)
		}
	}

	if eci.DescriptionOriginalLangage == language.Invalid {
		t.Warn("noDescription", "year", eci.Year, "nb", eci.Number)
	} else {
		for _, l := range t.Languages {
			if eci.Description[l] == nil {
				eci.Description[l] = eci.Description[eci.DescriptionOriginalLangage]
			}
		}
	}

	// Timeline
	for _, p := range dto.Progress {
		timeline := Timeline{
			Date:   p.Date.Time,
			Status: p.Status,
		}
		switch p.Status {
		case "REGISTERED":
			timeline.Register = registrationDoc
		case "CLOSED":
			if p.Note == "COLLECTION_EARLY_CLOSURE" {
				timeline.EarlyClose = true
			} else if p.Note != "" {
				t.Warn("unknwon.footnoteType", "year", info.year, "nb", info.number, "footnote", p.Note)
			}
		case "ANSWERED":
		case "COLLECTION_START_DATE", "INSUFFICIENT_SUPPORT", "ONGOING", "REJECTED", "SUBMITTED", "VERIFICATION", "WITHDRAWN":
			// ok
		default:
			t.Warn("unknwon.status", "year", info.year, "nb", info.number, "status", p.Status)
		}
		eci.Timeline = append(eci.Timeline, timeline)
	}
	slices.SortFunc(eci.Timeline, func(a, b Timeline) int {
		return cmp.Compare(a.Date.Unix(), b.Date.Unix())
	})

	// Set signature
	setSignature := func(entrys []signatureDTO) {
		for _, entry := range entrys {
			eci.Signature[entry.Country] = entry.Total
			eci.TotalSignature += entry.Total
		}
		date := time.Time{}
		for _, t := range eci.Timeline {
			if t.Status == "REGISTERED" {
				date = t.Date
			}
		}
		switch {
		case date_2024_07_06.Before(date):
			eci.Threshold = threshold_2024_07_06
		case date_2020_02_01.Before(date):
			eci.Threshold = threshold_2020_02_01
		case date_2020_01_01.Before(date):
			eci.Threshold = threshold_2020_01_01
		case date_2014_07_01.Before(date):
			eci.Threshold = threshold_2014_07_01
		case date_2012_04_01.Before(date):
			eci.Threshold = threshold_2012_04_01
		default:
			t.Warn("tooOldRegisterdate", "date", date, "year", eci.Year, "nb", eci.Number)
		}
		for c, sig := range eci.Signature {
			if eci.Threshold[c] <= sig {
				eci.ThresholdPassed++
			}
		}
	}
	if len(dto.Signatures.Entry) > 0 {
		eci.PaperSignaturesUpdate = dto.Signatures.UpdateDate.Time
		setSignature(dto.Signatures.Entry)
	} else {
		eci.ValidatedSignature = true
		setSignature(dto.Submission.Entry)
	}

	// Set image
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
		t.Warn("fetchImage", "err", "unknown format", "format", format, "logoID", logoID)
		return
	}

	eci.ImageWidth = strconv.Itoa(config.Width)
	eci.ImageHeight = strconv.Itoa(config.Height)
	eci.ImageData = data
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

type mimeTypeDTO string

func (m *mimeTypeDTO) UnmarshalText(data []byte) error {
	s := string(data)
	_, _, err := mime.ParseMediaType(s)
	if err != nil {
		return fmt.Errorf("mimeTypeDTO: %w", err)
	}
	*m = mimeTypeDTO(s)
	return nil
}
