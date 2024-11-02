package eu_ec_eci

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"maps"
	"mime"
	"net/url"
	"slices"
	"sniffle/common/resize0"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"sniffle/tool/sch"
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
	Threshold             *Threshold
	ThresholdRule         string
	ThresholdPass         [country.Len]bool
	ThresholdPassTotal    uint

	ImageName   string
	ImageWidth  string
	ImageHeight string
	ImageData   []byte

	ImageResizedName   string
	ImageResizedWidth  string
	ImageResizedHeight string
	ImageResizedData   []byte
}
type Description struct {
	Title       string
	PlainDesc   string
	SupportLink *url.URL
	Website     *url.URL
	FollowUp    *url.URL
	Objective   render.H
	AnnexDoc    *Document
	Annex       render.H
	Treaty      render.H
}
type Timeline struct {
	Date       time.Time
	Status     string
	EarlyClose bool
	Register   *[language.Len]*Document
	// Answer documents
	AnswerAnnex        *[language.Len]*Document
	AnswerResponse     *[language.Len]*Document
	AnswerPressRelease *[language.Len]*Document
}
type Threshold = [country.Len]uint
type Document struct {
	URL      *url.URL
	Language language.Language
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
	if tool.DevMode {
		j := tool.FetchAll(ctx, t, "", indexURL, nil, nil)
		t.WriteFile("/eu/ec/eci/src.json", j)

		var value any
		dec := json.NewDecoder(bytes.NewReader(j))
		dec.UseNumber()
		dec.Decode(&value)
		sch.Log(t.Logger.With("id", "index"), indexType, value)
	}
	if tool.FetchJSON(ctx, t, "", indexURL, nil, nil, &dto) {
		return nil
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
	dto := detailDTO{}
	fetchURL := fmt.Sprintf(detailURL, info.year, info.number)
	if tool.DevMode {
		j := tool.FetchAll(ctx, t, "", fetchURL, nil, nil)
		t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/src.json", info.year, info.number), j)
		var eciDTO any
		dec := json.NewDecoder(bytes.NewReader(j))
		dec.UseNumber()
		dec.Decode(&eciDTO)
		sch.Log(t.Logger.With("id", fmt.Sprintf("%d/%d", info.year, info.number)), eciType, eciDTO)
	}
	if tool.FetchJSON(ctx, t, "", fetchURL, nil, nil, &dto) {
		return nil
	}
	if tool.DevMode {
		dto.check(t.With("year", info.year, "nb", info.number))
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
	registrationDoc := new([language.Len]*Document)
	defaultAnnexDoc := (*Document)(nil)
	for _, desc := range dto.Description {
		if desc.SupportLink == desc.Website {
			desc.SupportLink = ""
		}
		supportLink := securehtml.ParseURL(desc.SupportLink)
		if supportLink != nil && supportLink.Host == "ec.europa.eu" {
			supportLink = nil
		}
		annexDoc := desc.AnnexDoc.Document(desc.Language)
		eci.Description[desc.Language] = &Description{
			Title:       desc.Title,
			PlainDesc:   securehtml.Text(desc.Objective, 200),
			SupportLink: supportLink,
			Website:     securehtml.ParseURL(desc.Website),
			Objective:   securehtml.Secure(desc.Objective),
			AnnexDoc:    annexDoc,
			Annex:       securehtml.Secure(desc.Annex),
			Treaty:      securehtml.TextWithURL(desc.Treaty),
		}
		if desc.Original {
			eci.DescriptionOriginalLangage = desc.Language
			defaultAnnexDoc = annexDoc
		}
		if u := securehtml.ParseURL(desc.Register.Url); u != nil {
			registrationDoc[desc.Language] = &Document{URL: u}
		} else if desc.Register.Document != nil {
			registrationDoc[desc.Language] = desc.Register.Document.Document(desc.Language)
		}
	}
	for _, desc := range eci.Description {
		if desc.AnnexDoc == nil {
			desc.AnnexDoc = defaultAnnexDoc
		}
	}

	for _, l := range t.Languages {
		if eci.Description[l] == nil {
			eci.Description[l] = eci.Description[eci.DescriptionOriginalLangage]
		}
		if registrationDoc[l] == nil {
			registrationDoc[l] = registrationDoc[eci.DescriptionOriginalLangage]
		}
	}

	answer := Timeline{}
	for _, link := range dto.Answer.Links {
		def := &Document{
			URL:      securehtml.ParseURL(link.DefaultLink),
			Language: link.DefaultLanguageCode,
		}

		translation := &[language.Len]*Document{}
		for l := range translation {
			translation[l] = def
		}
		for _, t := range link.Link {
			u := securehtml.ParseURL(t.Link)
			if u == nil {
				continue
			}
			translation[t.LanguageCode] = &Document{URL: u, Language: t.LanguageCode}
		}

		switch link.DefaultName {
		case "ANNEX":
			answer.AnswerAnnex = translation
		case "COMMUNICATION":
			answer.AnswerResponse = translation
		case "FOLLOW_UP":
			for l, desc := range eci.Description {
				u := translation[l].URL
				if u.Scheme == "https" && u.Host == "citizens-initiative.europa.eu" {
					u = u.JoinPath() // clone the url
					u.Path += "_" + l.String()
				}
				desc.FollowUp = u
			}
		case "PRESS_RELEASE":
			answer.AnswerPressRelease = translation
		}
	}

	// Timeline
	for _, p := range dto.Progress {
		if p.Date.Time.IsZero() {
			continue
		}
		timeline := Timeline{
			Date:   p.Date.Time,
			Status: p.Status,
		}
		switch timeline.Status {
		case "REGISTERED":
			timeline.Register = registrationDoc
		case "COLLECTION_START_DATE":
			timeline.Status = "ONGOING"
		case "CLOSED":
			if p.Note == "COLLECTION_EARLY_CLOSURE" {
				timeline.EarlyClose = true
			}
		case "ANSWERED":
			timeline.AnswerAnnex = answer.AnswerAnnex
			timeline.AnswerResponse = answer.AnswerResponse
			timeline.AnswerPressRelease = answer.AnswerPressRelease
		}
		eci.Timeline = append(eci.Timeline, timeline)
	}
	if !dto.Deadline.Time.IsZero() {
		eci.Timeline = append(eci.Timeline, Timeline{
			Date:   dto.Deadline.Time,
			Status: "DEADLINE",
		})
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
		registerDate := time.Time{}
		for _, t := range eci.Timeline {
			if t.Status == "REGISTERED" {
				registerDate = t.Date
			}
		}
		switch {
		case date_2024_07_06.Before(registerDate):
			eci.ThresholdRule = rule_since_2020_01_01
			eci.Threshold = &threshold_2024_07_06
		case date_2020_02_01.Before(registerDate):
			eci.ThresholdRule = rule_since_2020_01_01
			eci.Threshold = &threshold_2020_02_01
		case date_2020_01_01.Before(registerDate):
			eci.ThresholdRule = rule_since_2020_01_01
			eci.Threshold = &threshold_2020_01_01
		case date_2014_07_01.Before(registerDate):
			eci.ThresholdRule = rule_since_2012_04_01
			eci.Threshold = &threshold_2014_07_01
		case date_2012_04_01.Before(registerDate):
			eci.ThresholdRule = rule_since_2012_04_01
			eci.Threshold = &threshold_2012_04_01
		}
		for c, sig := range eci.Signature {
			if eci.Threshold[c] <= sig {
				eci.ThresholdPass[c] = true
				eci.ThresholdPassTotal++
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

	u := fmt.Sprintf(logoURL, logoID)
	data := tool.FetchAll(ctx, t, "", u, nil, nil)
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

	if resized := t.LongTask(resize0.Name, u, data); len(resized) != 0 {
		eci.ImageResizedName = "logo" + resize0.Extension
		eci.ImageResizedData = resized
		width, height := resize0.NewDimension(config.Width, config.Height)
		eci.ImageResizedWidth = strconv.Itoa(width)
		eci.ImageResizedHeight = strconv.Itoa(height)
	}
}

func (eci *ECIOut) countryByName(lang language.Language) []country.Country {
	name := translate.AllTranslation[lang].Country
	return slices.SortedFunc(maps.Keys(eci.Signature), func(a, b country.Country) int {
		return cmp.Compare(name[a], name[b])
	})
}

type dtoDate struct {
	Time time.Time
}

func (dto *dtoDate) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
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

type docDTO struct {
	Id       int
	Name     string
	Size     int
	MimeType mimeTypeDTO
}

func (doc *docDTO) Document(lang language.Language) *Document {
	if doc == nil {
		return nil
	}
	return &Document{
		URL: &url.URL{
			Scheme: "https",
			Host:   "register.eci.ec.europa.eu",
			Path:   "/core/api/register/document/" + strconv.Itoa(doc.Id),
		},
		Language: lang,
		Name:     doc.Name,
		MimeType: string(doc.MimeType),
		Size:     doc.Size,
	}
}

type mimeTypeDTO string

func (m *mimeTypeDTO) UnmarshalText(data []byte) error {
	s := string(data)
	mediatype, _, err := mime.ParseMediaType(s)
	if err != nil {
		return fmt.Errorf("mimeTypeDTO: %w", err)
	}
	*m = mimeTypeDTO(mediatype)
	return nil
}
