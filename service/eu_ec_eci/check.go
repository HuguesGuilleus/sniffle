package eu_ec_eci

import (
	"log/slog"
	"sniffle/common/country"
	"sniffle/front/translate"
	"sniffle/tool/language"
	"sniffle/tool/securehtml"
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
		ReplacedMember   []struct {
			FullName         string          `json:"fullName"`
			Type             string          `json:"type"`
			URL              string          `json:"email"`
			ResidenceCountry country.Country `json:"residenceCountry"`
			Start            dtoDate         `json:"startingDate"`
			End              dtoDate         `json:"endDate"`
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

// Report unknown value.
func (dto *detailDTO) check(logger *slog.Logger) {
	// Global Status
	if isUnknownStatus(dto.Status) {
		logger.Warn("check.unknwon", "status", dto.Status)
	}

	// Categorie
	knownCategorie := translate.T[language.English].EU_EC_ECI.Categorie
	categories := make(map[string]bool, len(dto.Categories))
	for _, c := range dto.Categories {
		if knownCategorie[c.CategoryType] == "" {
			logger.Warn("check.unknwon", "categorie", c.CategoryType)
		}
		if categories[c.CategoryType] {
			logger.Warn("check.multiple", "categorie", c.CategoryType)
		}
		categories[c.CategoryType] = true
	}

	// Description: Language + Original + has registration
	originalDescription := language.Invalid
	for i, desc := range dto.Description {
		if desc.Language == language.Invalid {
			logger.Warn("check.wrongLangage", "i", i)
		}

		if desc.Original {
			if originalDescription != language.Invalid {
				logger.Warn("check.multipleOriginalDescription",
					"first", originalDescription,
					"lang", desc.Language)
			} else {
				originalDescription = desc.Language
			}
		}

		if desc.Title == "" {
			logger.Warn("check.noTitle", "i", i, "lang", desc.Language)
		}
		if desc.Objective == "" {
			logger.Warn("check.noObjective", "i", i, "lang", desc.Language)
		}

		noRegistration := true
		if u := desc.Register.Url; u != "" && securehtml.ParseURL(u) == nil {
			logger.Warn("check.wrongURL(description[$].Register.Url)",
				"i", i, "lang", desc.Language,
				"url", u)
		} else {
			noRegistration = false
		}
		if doc := desc.Register.Document; doc != nil && doc.Id == 0 {
			logger.Warn("check.registrationDocID", "i", i, "lang", desc.Language, "id", doc.Id)
		}
		if noRegistration {
			logger.Warn("check.noRegistration", "i", i, "lang", desc.Language)
		} else if desc.Register.Url != "" && false {
			logger.Warn("check.multipleRegistration", "i", i, "lang", desc.Language)
		}
	}
	if originalDescription == language.Invalid {
		logger.Warn("check.noOriginalDescription")
	}

	// Timeline: Status + Note
	status := make(map[string]bool, len(dto.Progress))
	for _, p := range dto.Progress {
		if isUnknownStatus(p.Status) {
			logger.Warn("check.unknown", "status", p.Status)
		}
		if status[p.Status] {
			logger.Warn("check.multiple", "status", p.Status)
		}
		status[p.Status] = true

		if p.Status == "CLOSED" && p.Note == "COLLECTION_EARLY_CLOSURE" {
			// ok
		} else if p.Note != "" {
			logger.Warn("check.unknown", "with-status", p.Status, "note", p.Note)
		}

		if p.Date.Time.Before(date_2012_04_01) &&
			!p.Date.Time.IsZero() &&
			p.Status != "COLLECTION_START_DATE" {
			logger.Warn("check.dateBefore_2012_04_01", "with-status", p.Status, "date", p.Date.Time)
		}
	}
}

func isUnknownStatus(status string) bool {
	switch status {
	case "REGISTERED", "COLLECTION_START_DATE", "CLOSED", "ANSWERED", "INSUFFICIENT_SUPPORT", "ONGOING", "REJECTED", "SUBMITTED", "VERIFICATION", "WITHDRAWN":
		return false
	default:
		return true
	}
}
