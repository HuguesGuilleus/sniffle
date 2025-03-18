package eu_ec_eci

import (
	"cmp"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"net/url"
	"slices"
	"sniffle/common"
	"sniffle/common/country"
	"sniffle/common/language"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/render"
	"sniffle/tool/sch"
	"sniffle/tool/securehtml"
	"strconv"
	"strings"
	"time"
)

const (
	detailURL = "https://register.eci.ec.europa.eu/core/api/register/details/%d/%06d"
	logoURL   = "https://register.eci.ec.europa.eu/core/api/register/logo/%d"
)

type ECIOut struct {
	// Identifier
	Year   int
	Number int
	// Last update of information
	LastUpdate time.Time
	// Current status
	Status string
	// A sorted and uniq slice of categories
	Categorie []string
	// The image
	Image *common.ResizedImage

	// Description text in all language
	Description map[language.Language]*Description
	// The based language used to write the ECI text.
	OriginalLangage language.Language

	Timeline []Event

	Signature []Signature
	// Total of validated signature
	TotalSignature uint
	// Date of the last organisators paper signatures update.
	// Can be zero
	PaperSignaturesUpdate time.Time

	ThresholdRule      string
	ThresholdPassTotal uint

	Members []Member

	FundingTotal    float64
	FundingUpdate   time.Time
	FundingDocument *Document
	Sponsor         []Sponsor
}
type Description struct {
	Title       string
	SupportLink *url.URL
	Website     *url.URL
	FollowUp    *url.URL

	PlainDesc string
	Objective render.H
	Annex     render.H
	Treaty    render.H

	AnnexDoc   *Document
	DraftLegal *Document
}
type Event struct {
	Status string
	Date   time.Time

	// The registration document
	Register *[language.Len]*Document

	// EarlyClose when organisator close signatures reception
	// If the status is CLOSED
	EarlyClose bool

	// Answer documents
	AnswerAnnex        *[language.Len]*Document
	AnswerResponse     *[language.Len]*Document
	AnswerPressRelease *[language.Len]*Document
}
type Signature struct {
	Country country.Country
	// Number of signature
	Count uint
	// After the deadline
	After bool
	// Threashold for this country at registration date.
	Threshold uint
	// If Count >= Threshold
	ThresholdPass bool
}
type Member struct {
	// "MEMBER" | "SUBSTITUTE" | "REPRESENTATIVE" | "OTHER" | "DPO" | "LEGAL_ENTITY"
	Type string

	FullName string

	// Nothing, HTTP.S or mailto:...
	// Href used in anchor href attribute
	HrefURL string
	// DisplayURL is used to print.
	DisplayURL string

	// Maybe zero.
	Start time.Time
	End   time.Time

	ResidenceCountry country.Country

	// Privacy applied
	Privacy bool

	// Only one depth level.
	Replaced *Member
}
type Sponsor struct {
	// Name of the sponsor.
	// If anonyous, the name is empty.
	Name string
	// IsPrivate or is an organisasion
	IsPrivate bool
	Amount    float64
	Date      time.Time
}
type Document struct {
	URL      *url.URL
	Language language.Language
	Name     string
	Size     int
	MimeType string
}

func fetchDetail(t *tool.Tool, info indexItem) *ECIOut {
	request := fetch.URL(fmt.Sprintf(detailURL, info.year, info.number))
	if tool.DevMode {
		t.WriteFile(
			fmt.Sprintf("/eu/ec/eci/%d/%d/src.json", info.year, info.number),
			tool.FetchAll(t, request),
		)
	}
	dto := acceptedDTO{}
	if tool.FetchJSON(t, eciType, &dto, request) {
		return nil
	}

	eci := &ECIOut{
		Year:       info.year,
		Number:     info.number,
		LastUpdate: dto.LastUpdate.Time,
		Status:     dto.Status,
		Categorie:  make([]string, len(dto.Categories)),
		Image:      fetchImage(t, dto.Logo.ID),
	}

	// Categorie
	for i, entry := range dto.Categories {
		eci.Categorie[i] = entry.CategoryType
	}
	slices.Sort(eci.Categorie)

	// Description
	eci.Description = make(map[language.Language]*Description)
	for _, desc := range dto.Description {
		if desc.Original {
			eci.OriginalLangage = desc.Language
		}
		if desc.SupportLink == desc.Website {
			desc.SupportLink = ""
		}
		supportLink := securehtml.ParseURL(desc.SupportLink)
		if eci.Status != "ONGOING" || (supportLink != nil && supportLink.Host == "ec.europa.eu") {
			supportLink = nil
		}
		eci.Description[desc.Language] = &Description{
			Title:       desc.Title,
			SupportLink: supportLink,
			Website:     securehtml.ParseURL(desc.Website),

			PlainDesc: securehtml.Text(desc.Objective, 200),
			Objective: securehtml.Secure(desc.Objective),
			Annex:     securehtml.Secure(desc.Annex),
			Treaty:    securehtml.TextWithURL(desc.Treaty),

			AnnexDoc:   desc.AnnexDoc.Document(desc.Language),
			DraftLegal: desc.DraftLegal.Document(desc.Language),
		}
	}
	defaultDesc := eci.Description[eci.OriginalLangage]
	for _, desc := range eci.Description {
		if desc.AnnexDoc == nil {
			desc.AnnexDoc = defaultDesc.AnnexDoc
		}
		if desc.DraftLegal == nil {
			desc.DraftLegal = defaultDesc.DraftLegal
		}
	}

	// if pre answer, get the follow up link
	if du := securehtml.ParseURL(dto.PreAnswer.Links[0].DefaultLink); du != nil {
		for l, desc := range eci.Description {
			if desc != nil {
				u := (*du)
				u.Path += "_" + l.String()
				desc.FollowUp = &u
			}
		}
	}

	// Registration documents
	registrationDoc := new([language.Len]*Document)
	for _, desc := range dto.Description {
		if u := securehtml.ParseURL(desc.Register.URL); u != nil {
			registrationDoc[desc.Language] = &Document{URL: u}
		} else {
			registrationDoc[desc.Language] = desc.Register.Document.Document(desc.Language)
		}
	}

	// Get answer documents
	answer := Event{} // wraper for all documents
	for _, link := range dto.Answer.Links {
		translation := &[language.Len]*Document{}
		def := &Document{
			URL:      securehtml.ParseURL(link.DefaultLink),
			Language: link.DefaultLanguageCode,
		}
		for l := range translation {
			translation[l] = def
		}
		for _, t := range link.Link {
			u := securehtml.ParseURL(t.Link)
			if u == nil {
				continue
			}
			translation[t.Language] = &Document{URL: u, Language: t.Language}
		}

		switch link.Kind {
		case "ANNEX":
			answer.AnswerAnnex = translation
		case "COMMUNICATION":
			answer.AnswerResponse = translation
		case "PRESS_RELEASE":
			answer.AnswerPressRelease = translation
		case "FOLLOW_UP":
			for l, desc := range eci.Description {
				u := translation[l].URL
				if u.Scheme == "https" && u.Host == "citizens-initiative.europa.eu" {
					u = &url.URL{
						Scheme: "https",
						Host:   "citizens-initiative.europa.eu",
						Path:   u.Path + "_" + l.String(),
					}
				}
				desc.FollowUp = u
			}
		}
	}

	// Timeline
	registerDate := time.Time{} // use for Threshold
	for _, p := range dto.Progress {
		if p.Date.Time.IsZero() {
			continue
		}
		e := Event{
			Date:   p.Date.Time,
			Status: p.Status,
		}
		switch e.Status {
		case "REGISTERED":
			registerDate = e.Date
			e.Register = registrationDoc
		case "CLOSED":
			e.EarlyClose = p.Note == "COLLECTION_EARLY_CLOSURE"
		case "ANSWERED":
			e.AnswerAnnex = answer.AnswerAnnex
			e.AnswerResponse = answer.AnswerResponse
			e.AnswerPressRelease = answer.AnswerPressRelease
		}
		eci.Timeline = append(eci.Timeline, e)
	}
	if !dto.Deadline.Time.IsZero() {
		eci.Timeline = append(eci.Timeline, Event{
			Date:   dto.Deadline.Time,
			Status: "DEADLINE",
		})
	}
	slices.SortFunc(eci.Timeline, func(a, b Event) int {
		return cmp.Compare(a.Date.Unix(), b.Date.Unix())
	})

	// Set signature
	if dto.SosReport != nil {
		eci.setSignatures(dto.SosReport, registerDate)
	} else if dto.Submission != nil {
		eci.setSignatures(dto.Submission, registerDate)
	}

	// Members
	eci.Members = make([]Member, len(dto.Members))
	for i, entry := range dto.Members {
		replacedMember := (*Member)(nil)
		if len(entry.ReplacedMember) == 1 {
			r := entry.ReplacedMember[0]
			hrefURL, displayURL := memberURL(r.URL)
			replacedMember = &Member{
				FullName:         strings.TrimSpace(r.FullName),
				Type:             r.Type,
				HrefURL:          hrefURL,
				DisplayURL:       displayURL,
				ResidenceCountry: r.ResidenceCountry,
				Start:            r.Start.Time,
				End:              r.End.Time,
				Privacy:          r.Privacy,
			}
		}
		hrefURL, displayURL := memberURL(entry.URL)
		eci.Members[i] = Member{
			FullName:         strings.TrimSpace(entry.FullName),
			Type:             entry.Type,
			HrefURL:          hrefURL,
			DisplayURL:       displayURL,
			ResidenceCountry: entry.ResidenceCountry,
			Start:            entry.Start.Time,
			Privacy:          entry.Privacy,
			Replaced:         replacedMember,
		}
	}

	// Funding
	if fund := dto.Funding; !fund.LastUpdate.Time.IsZero() {
		eci.FundingUpdate = fund.LastUpdate.Time
		eci.FundingTotal = fund.TotalAmount
		eci.FundingDocument = fund.Document.Document(0)
		eci.Sponsor = make([]Sponsor, len(fund.Sponsors))
		for i, s := range fund.Sponsors {
			name := s.Name
			if s.Anonymized {
				name = ""
			}
			eci.Sponsor[i] = Sponsor{
				Name:      name,
				IsPrivate: s.PrivateSponsor,
				Amount:    s.Amount,
				Date:      s.Date.Time,
			}
		}
	}

	return eci
}

func fetchImage(t *tool.Tool, logoID int) *common.ResizedImage {
	if logoID == 0 {
		return nil
	}
	return common.FetchImage(t, fetch.URL(fmt.Sprintf(logoURL, logoID)))
}

func (eci *ECIOut) setSignatures(dto *signatureDTO, registerDate time.Time) {
	eci.TotalSignature = dto.TotalSignatures
	eci.PaperSignaturesUpdate = dto.PaperSignaturesUpdate.Time

	threshold := [country.Len]uint{}
	switch {
	case date_2024_07_06.Before(registerDate):
		eci.ThresholdRule = rule_since_2020_01_01
		threshold = threshold_2024_07_06
	case date_2020_02_01.Before(registerDate):
		eci.ThresholdRule = rule_since_2020_01_01
		threshold = threshold_2020_02_01
	case date_2020_01_01.Before(registerDate):
		eci.ThresholdRule = rule_since_2020_01_01
		threshold = threshold_2020_01_01
	case date_2014_07_01.Before(registerDate):
		eci.ThresholdRule = rule_since_2012_04_01
		threshold = threshold_2014_07_01
	case date_2012_04_01.Before(registerDate):
		eci.ThresholdRule = rule_since_2012_04_01
		threshold = threshold_2012_04_01
	}

	eci.Signature = make([]Signature, len(dto.Entry))
	for i, entry := range dto.Entry {
		sig := Signature{
			Country:   entry.Country,
			Count:     entry.Count,
			After:     entry.After,
			Threshold: threshold[entry.Country],
		}
		if !sig.After && sig.Count >= sig.Threshold {
			sig.ThresholdPass = true
			eci.ThresholdPassTotal++
		}
		eci.Signature[i] = sig
	}
	slices.SortFunc(eci.Signature, func(a, b Signature) int {
		return cmp.Compare(a.Country, b.Country)
	})
}

func memberURL(s string) (href, display string) {
	if s == "email@anonymised" {
		return "", ""
	} else if s != "" && sch.AnyMail().Match(s) == nil {
		return "mailto:" + s, s
	} else if u := securehtml.ParseURL(s); u != nil {
		s := u.String()
		return s, s
	}
	return "", ""
}

// All available langs for this ICE and include in [translate.Langs]
func (eci *ECIOut) Langs() []language.Language {
	langs := make([]language.Language, 0, len(translate.Langs))
	for _, l := range translate.Langs {
		if eci.Description[l] != nil {
			langs = append(langs, l)
		}
	}
	return langs
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
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Size     int    `json:"size"`
	MimeType string `json:"mimeType"`
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
