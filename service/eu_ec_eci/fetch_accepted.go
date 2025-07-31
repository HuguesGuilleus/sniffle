package eu_ec_eci

import (
	"cmp"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/HuguesGuilleus/sniffle/common/country"
	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/common/rimage"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fetch"
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/HuguesGuilleus/sniffle/tool/sch"
	"github.com/HuguesGuilleus/sniffle/tool/securehtml"
)

const (
	detailURL    = "https://register.eci.ec.europa.eu/core/api/register/details/%d/%06d"
	logoURL      = "https://register.eci.ec.europa.eu/core/api/register/logo/%d"
	plainDescMax = 200
)

func fetchAllAcepted(t *tool.Tool) map[uint][]*ECIOut {
	infoIndex := fetchAcceptedIndex(t)
	wg := sync.WaitGroup{}
	wg.Add(1 + len(infoIndex))
	go func() {
		defer wg.Done()
		checkThreashold(t)
	}()
	eciByYear := make(map[uint][]*ECIOut)
	mutex := sync.Mutex{}
	for _, info := range infoIndex {
		go func(info indexItem) {
			defer wg.Done()
			eci := fetchDetail(t, info)
			if eci == nil {
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			eciByYear[eci.Year] = append(eciByYear[eci.Year], eci)
		}(info)
	}
	wg.Wait()

	for _, byYear := range eciByYear {
		slices.SortFunc(byYear, func(a, b *ECIOut) int {
			return cmp.Compare(b.Number, a.Number)
		})
	}

	return eciByYear
}

type ECIOut struct {
	// Identifier
	ID     uint
	Year   uint
	Number uint
	// Last update of information
	LastUpdate time.Time
	// Current status
	Status string
	// A sorted and uniq slice of categories
	Categorie []string
	// The image
	Image *rimage.Image

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
	// Threshold information
	Threshold          *Threshold
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
	Register            *[language.Len]*Document
	RegisterCorrigendum *[language.Len]*Document

	// EarlyClose when organisator close signatures reception
	// If the status is CLOSED
	EarlyClose bool
	// Because of COVID, some ICE get collecting signature extra delay.
	// Legal text source.
	// If nil, no extra delay.
	ExtraDelay []component.Legal

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

type acceptedDTO struct {
	Status     string  `json:"status"`
	LastUpdate dtoTime `json:"latestUpdateDate"`
	Deadline   dtoDate `json:"deadline"`
	Logo       struct {
		ID int `json:"id"`
	} `json:"logo"`
	Categories []struct {
		CategoryType string `json:"categoryType"`
	} `json:"categories"`

	Description []struct {
		Original    bool              `json:"Original"`
		Language    language.Language `json:"languageCode"`
		Title       string            `json:"title"`
		SupportLink string            `json:"supportLink"`
		Website     string            `json:"website"`
		Objective   string            `json:"objectives"`
		AnnexDoc    *docDTO           `json:"additionalDocument"`
		DraftLegal  *docDTO           `json:"draftLegal"`
		Annex       string            `json:"annexText"`
		Treaty      string            `json:"treaties"`
		Register    struct {
			CELEX       string  `json:"celex"`
			Corrigendum string  `json:"corrigendum"`
			Document    *docDTO `json:"document"`
		} `json:"commissionDecision"`
	} `json:"linguisticVersions"`
	PreAnswer struct {
		Links [1]struct {
			DefaultLink string `json:"defaultLink"`
		} `json:"links"`
	} `json:"preAnswer"`

	Progress []struct {
		Status string  `json:"Name"`
		Note   string  `json:"footnoteType"`
		Date   dtoDate `json:"date"`
	} `json:"progress"`
	Answer struct {
		Links []struct {
			DefaultLanguageCode language.Language `json:"defaultLanguageCode"`
			Kind                string            `json:"defaultName"`
			DefaultLink         string            `json:"defaultLink"`
			Link                []struct {
				Language language.Language `json:"languageCode"`
				Link     string            `json:"link"`
			} `json:"link"`
		} `json:"links"`
	} `json:"answer"`

	SosReport  *signatureDTO `json:"sosReport"`
	Submission *signatureDTO `json:"submission"`

	Members []struct {
		FullName         string          `json:"fullName"`
		Type             string          `json:"type"`
		URL              string          `json:"email"`
		ResidenceCountry country.Country `json:"ResidenceCountry"`
		Start            dtoDate         `json:"startingDate"`
		Privacy          bool            `json:"privacyApplied"`
		ReplacedMember   []struct {
			FullName         string          `json:"fullName"`
			Type             string          `json:"type"`
			URL              string          `json:"email"`
			ResidenceCountry country.Country `json:"residenceCountry"`
			Start            dtoDate         `json:"startingDate"`
			End              dtoDate         `json:"endDate"`
			Privacy          bool            `json:"privacyApplied"`
		} `json:"replacedMember"`
	} `json:"members"`

	Funding struct {
		LastUpdate dtoDate `json:"lastUpdate"`
		Sponsors   []struct {
			Amount         float64 `json:"amount"`
			Date           dtoDate `json:"date"`
			Name           string  `json:"name"`
			PrivateSponsor bool    `json:"privateSponsor"`
			Anonymized     bool    `json:"anonymized"`
		}
		TotalAmount float64 `json:"totalAmount"`
		Document    *docDTO `json:"document"`
	} `json:"funding"`
}
type signatureDTO struct {
	PaperSignaturesUpdate dtoDate `json:"updateDate"`
	TotalSignatures       uint    `json:"totalSignatures"`

	Entry []struct {
		Country country.Country `json:"countryCodeType"`
		Count   uint            `json:"total"`
		After   bool            `json:"afterSubmission"`
	}
}

func fetchDetail(t *tool.Tool, info indexItem) *ECIOut {
	request := fetch.Fmt(detailURL, info.year, info.number)
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
		ID:         info.id,
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

			PlainDesc: securehtml.Text(desc.Objective, plainDescMax),
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
	registrationDocCorrigendum := new([language.Len]*Document)
	for _, desc := range dto.Description {
		if desc.Register.Document != nil {
			registrationDoc[desc.Language] = desc.Register.Document.Document(desc.Language)
		} else {
			registrationDoc[desc.Language] = docFromCelex(desc.Register.CELEX, desc.Language)
			if desc.Register.Corrigendum != "" {
				registrationDocCorrigendum[desc.Language] = docFromCelex(desc.Register.Corrigendum, desc.Language)
			}
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
			e.RegisterCorrigendum = registrationDocCorrigendum
		case "CLOSED":
			e.EarlyClose = p.Note == "COLLECTION_EARLY_CLOSURE"
			e.ExtraDelay = extraDelayMap[info.id]
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

func fetchImage(t *tool.Tool, logoID int) *rimage.Image {
	if logoID == 0 {
		return nil
	}
	return rimage.New(t, fmt.Sprintf(logoURL, logoID))
}

func (eci *ECIOut) setSignatures(dto *signatureDTO, registerDate time.Time) {
	eci.TotalSignature = dto.TotalSignatures
	eci.PaperSignaturesUpdate = dto.PaperSignaturesUpdate.Time

	for _, t := range thresholds {
		if t.Begin.Before(registerDate) {
			eci.Threshold = t
			break
		}
	}

	eci.Signature = make([]Signature, len(dto.Entry))
	for i, entry := range dto.Entry {
		sig := Signature{
			Country:   entry.Country,
			Count:     entry.Count,
			After:     entry.After,
			Threshold: eci.Threshold.Data[entry.Country],
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

func docFromCelex(celex string, l language.Language) *Document {
	return &Document{URL: &url.URL{
		Scheme: "https",
		Host:   "eur-lex.europa.eu",
		Path:   "/legal-content/" + l.Upper() + "/TXT/",
		RawQuery: url.Values{
			"uri": {"CELEX:" + celex},
		}.Encode(),
	}}
}

func memberURL(s string) (href, display string) {
	if s == "email@anonymised" || s == "anonymised@anonymised" {
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
