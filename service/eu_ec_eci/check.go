package eu_ec_eci

import (
	"sniffle/common/country"
	"sniffle/common/language"
)

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
			URL         string  `json:"url"`
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

	Signatures struct {
		UpdateDate dtoDate
		Entry      []signatureDTO
	} `json:"sosReport"`
	Submission struct {
		Entry []signatureDTO
	}

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
	Country country.Country `json:"countryCodeType"`
	Total   uint
}
