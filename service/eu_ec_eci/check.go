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
	}
}
type signatureDTO struct {
	Country country.Country `json:"countryCodeType"`
	Total   uint
}
