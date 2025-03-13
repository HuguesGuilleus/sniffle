package eu_ec_eci

import (
	"sniffle/common/country"
	"sniffle/tool/language"
)

type detailDTO struct {
	Status      string
	LastUpdate  dtoTime `json:"latestUpdateDate"`
	Deadline    dtoDate
	Categories  []struct{ CategoryType string }
	Description []struct {
		Original    bool
		Language    language.Language `json:"languageCode"`
		Title       string
		SupportLink string
		Website     string
		Objective   string  `json:"objectives"`
		AnnexDoc    *docDTO `json:"additionalDocument"`
		DraftLegal  *docDTO `json:"draftLegal"`
		Annex       string  `json:"annexText"`
		Treaty      string  `json:"treaties"`
		Register    struct {
			Url      string
			Document *docDTO
		} `json:"commissionDecision"`
	} `json:"linguisticVersions"`
	Members []struct {
		FullName         string
		Type             string
		URL              string `json:"email"`
		ResidenceCountry country.Country
		Start            dtoDate `json:"startingDate"`
		Privacy          bool    `json:"privacyApplied"`
		ReplacedMember   []struct {
			FullName         string          `json:"fullName"`
			Type             string          `json:"type"`
			URL              string          `json:"email"`
			ResidenceCountry country.Country `json:"residenceCountry"`
			Start            dtoDate         `json:"startingDate"`
			End              dtoDate         `json:"endDate"`
			Privacy          bool            `json:"privacyApplied"`
		}
	}
	Progress []struct {
		Status string `json:"Name"`
		Note   string `json:"footnoteType"`
		Date   dtoDate
	}
	Signatures struct {
		UpdateDate dtoDate
		Entry      []signatureDTO
	} `json:"sosReport"`
	Submission struct {
		Entry []signatureDTO
	}
	Funding struct {
		LastUpdate dtoDate
		Sponsors   []struct {
			Amount         float64
			Date           dtoDate
			Name           string
			PrivateSponsor bool
			Anonymized     bool
		}
		TotalAmount float64
		Document    *docDTO
	}
	Answer struct {
		Links []struct {
			DefaultLanguageCode language.Language
			DefaultName         string
			DefaultLink         string
			Link                []struct {
				LanguageCode language.Language
				Link         string
			}
		}
	}
}
type signatureDTO struct {
	Country country.Country `json:"countryCodeType"`
	Total   uint
}
