package eu_ec_ice

import (
	"context"
	"net/url"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/testingtool"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var fetcher = testingtool.TestFetcher{
	indexURL: []byte(`{
		"entries": [
			{
				"year": "2024",
				"number": "000008"
			},
			{
				"year": "2024",
				"number": "000009",
				"logo": { "id": 8846 }
			}
		]
	}`),
	"https://register.eci.ec.europa.eu/core/api/register/logo/8846": []byte(`image8846`),
	"https://register.eci.ec.europa.eu/core/api/register/details/2024/000009": []byte(`{
		"status": "ONGOING",
		"latestUpdateDate": "24/07/2024 13:52",
		"linguisticVersions": [
			{
				"original": true,
				"languageCode": "EN",
				"title": "Title",
				"objectives": "<p>objectives</p>",
				"annexText": "<ul><li>arg 1: BECAUSE!!!</li><li>arg 2 ...</li></ul>",
				"treaties": "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
				"website": "https://furfreeeurope.eu/"
			},
			{
				"original": false,
				"languageCode": "FR",
				"title": "Titre",
				"objectives": "<p>Objectifs</p>",
				"annexText": "<ul><li>arg 1: PARCE QUE!!!</li><li>arg 2 ...</li></ul>",
				"treaties": "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE"
			}
		],
		"categories": [
			{ "categoryType": "SANTE" },
			{ "categoryType": "TRADE" },
			{ "categoryType": "SANTE" },
			{ "categoryType": "AGRI" }
		],
		"sosReport": {
			"entry": [
				{ "countryCodeType": "FI", "total": 2171 },
				{ "countryCodeType": "RO", "total": 1362 },
				{ "countryCodeType": "CY", "total": 60 },
				{ "countryCodeType": "LU", "total": 354 },
				{ "countryCodeType": "PL", "total": 1879 },
				{ "countryCodeType": "IT", "total": 4940 },
				{ "countryCodeType": "HR", "total": 628 },
				{ "countryCodeType": "SE", "total": 1642 },
				{ "countryCodeType": "FR", "total": 24384 },
				{ "countryCodeType": "BG", "total": 196 },
				{ "countryCodeType": "HU", "total": 630 },
				{ "countryCodeType": "DK", "total": 866 },
				{ "countryCodeType": "IE", "total": 1757 },
				{ "countryCodeType": "GR", "total": 551 },
				{ "countryCodeType": "NL", "total": 3856 },
				{ "countryCodeType": "BE", "total": 2101 },
				{ "countryCodeType": "ES", "total": 11608 },
				{ "countryCodeType": "SK", "total": 434 },
				{ "countryCodeType": "DE", "total": 12985 },
				{ "countryCodeType": "LT", "total": 427 },
				{ "countryCodeType": "CZ", "total": 367 },
				{ "countryCodeType": "LV", "total": 149 },
				{ "countryCodeType": "SI", "total": 775 },
				{ "countryCodeType": "AT", "total": 1099 },
				{ "countryCodeType": "PT", "total": 650 },
				{ "countryCodeType": "EE", "total": 244 },
				{ "countryCodeType": "MT", "total": 61 }
			]
		}
	}`),
}

func TestFetchIndex(t *testing.T) {
	items, err := fetchIndex(context.Background(), fetcher)
	assert.NoError(t, err)
	assert.Equal(t, []indexItem{
		{year: 2024, number: 8},
		{year: 2024, number: 9, logoID: 8846},
	}, items)
}

func TestFetchDetail(t *testing.T) {
	ice, err := fetchDetail(context.Background(), fetcher, indexItem{
		year:   2024,
		number: 9,
		logoID: 8846,
	})
	assert.NoError(t, err)
	assert.Equal(t, &ICEOut{
		Year:       2024,
		Number:     9,
		LastUpdate: time.Date(2024, time.July, 24, 13, 52, 0, 0, time.UTC),
		Status:     "ONGOING",
		Categorie:  []string{"AGRI", "SANTE", "TRADE"},

		Description: map[language.Langage]*Description{
			language.English: {
				Title: "Title",
				Website: &url.URL{
					Scheme: "https",
					Host:   "furfreeeurope.eu",
					Path:   "/",
				},
				Objective: "<p>objectives</p>",
				Annex:     "<ul><li>arg 1: BECAUSE!!!</li><li>arg 2 ...</li></ul>",
				Treaty:    "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
			},
			language.French: {
				Title:     "Titre",
				Objective: "<p>Objectifs</p>",
				Annex:     "<ul><li>arg 1: PARCE QUE!!!</li><li>arg 2 ...</li></ul>",
				Treaty:    "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
			},
		},
		DescriptionOriginalLangage: language.English,

		TotalSignature: 76176,
		Signature: map[country.Country]uint{
			country.Finland:     2171,
			country.Romania:     1362,
			country.Cyprus:      60,
			country.Luxembourg:  354,
			country.Poland:      1879,
			country.Italy:       4940,
			country.Croatia:     628,
			country.Sweden:      1642,
			country.France:      24384,
			country.Bulgaria:    196,
			country.Hungary:     630,
			country.Denmark:     866,
			country.Ireland:     1757,
			country.Greece:      551,
			country.Netherlands: 3856,
			country.Belgium:     2101,
			country.Spain:       11608,
			country.Slovakia:    434,
			country.Germany:     12985,
			country.Lithuania:   427,
			country.Czechia:     367,
			country.Latvia:      149,
			country.Slovenia:    775,
			country.Austria:     1099,
			country.Portugal:    650,
			country.Estonia:     244,
			country.Malta:       61,
		},

		Image: []byte(`image8846`),
	}, ice)
	assert.Same(t, ice.Description[language.English], ice.GetOriginalDescription())
}
