package eu_ec_eci

import (
	"context"
	"net/url"
	"sniffle/tool"
	"sniffle/tool/country"
	"sniffle/tool/fetch"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var image3x1PNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x01, 0x01, 0x03, 0x00, 0x00, 0x00, 0x21, 0x2e, 0x86,
	0xf7, 0x00, 0x00, 0x00, 0x04, 0x67, 0x41, 0x4d, 0x41, 0x00, 0x00, 0xb1, 0x8f, 0x0b, 0xfc, 0x61,
	0x05, 0x00, 0x00, 0x00, 0x20, 0x63, 0x48, 0x52, 0x4d, 0x00, 0x00, 0x7a, 0x26, 0x00, 0x00, 0x80,
	0x84, 0x00, 0x00, 0xfa, 0x00, 0x00, 0x00, 0x80, 0xe8, 0x00, 0x00, 0x75, 0x30, 0x00, 0x00, 0xea,
	0x60, 0x00, 0x00, 0x3a, 0x98, 0x00, 0x00, 0x17, 0x70, 0x9c, 0xba, 0x51, 0x3c, 0x00, 0x00, 0x00,
	0x06, 0x50, 0x4c, 0x54, 0x45, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0x41, 0x1d, 0x34, 0x11, 0x00,
	0x00, 0x00, 0x01, 0x62, 0x4b, 0x47, 0x44, 0x01, 0xff, 0x02, 0x2d, 0xde, 0x00, 0x00, 0x00, 0x07,
	0x74, 0x49, 0x4d, 0x45, 0x07, 0xe8, 0x08, 0x0a, 0x0a, 0x26, 0x23, 0x05, 0x20, 0xef, 0x6a, 0x00,
	0x00, 0x00, 0x0a, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0x60, 0x00, 0x00, 0x00, 0x02, 0x00,
	0x01, 0xe2, 0x21, 0xbc, 0x33, 0x00, 0x00, 0x00, 0x25, 0x74, 0x45, 0x58, 0x74, 0x64, 0x61, 0x74,
	0x65, 0x3a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x00, 0x32, 0x30, 0x32, 0x34, 0x2d, 0x30, 0x38,
	0x2d, 0x31, 0x30, 0x54, 0x31, 0x30, 0x3a, 0x33, 0x38, 0x3a, 0x32, 0x31, 0x2b, 0x30, 0x30, 0x3a,
	0x30, 0x30, 0x2d, 0x52, 0xe8, 0x7e, 0x00, 0x00, 0x00, 0x25, 0x74, 0x45, 0x58, 0x74, 0x64, 0x61,
	0x74, 0x65, 0x3a, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x79, 0x00, 0x32, 0x30, 0x32, 0x34, 0x2d, 0x30,
	0x38, 0x2d, 0x31, 0x30, 0x54, 0x31, 0x30, 0x3a, 0x33, 0x38, 0x3a, 0x32, 0x31, 0x2b, 0x30, 0x30,
	0x3a, 0x30, 0x30, 0x5c, 0x0f, 0x50, 0xc2, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

var image3x1JPG = []byte{
	0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01,
	0x00, 0x01, 0x00, 0x00, 0xff, 0xdb, 0x00, 0x43, 0x00, 0x03, 0x02, 0x02, 0x02, 0x02, 0x02, 0x03,
	0x02, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03, 0x04, 0x06, 0x04, 0x04, 0x04, 0x04, 0x04, 0x08, 0x06,
	0x06, 0x05, 0x06, 0x09, 0x08, 0x0a, 0x0a, 0x09, 0x08, 0x09, 0x09, 0x0a, 0x0c, 0x0f, 0x0c, 0x0a,
	0x0b, 0x0e, 0x0b, 0x09, 0x09, 0x0d, 0x11, 0x0d, 0x0e, 0x0f, 0x10, 0x10, 0x11, 0x10, 0x0a, 0x0c,
	0x12, 0x13, 0x12, 0x10, 0x13, 0x0f, 0x10, 0x10, 0x10, 0xff, 0xdb, 0x00, 0x43, 0x01, 0x03, 0x03,
	0x03, 0x04, 0x03, 0x04, 0x08, 0x04, 0x04, 0x08, 0x10, 0x0b, 0x09, 0x0b, 0x10, 0x10, 0x10, 0x10,
	0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10,
	0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10,
	0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0xff, 0xc0,
	0x00, 0x11, 0x08, 0x00, 0x01, 0x00, 0x03, 0x03, 0x01, 0x11, 0x00, 0x02, 0x11, 0x01, 0x03, 0x11,
	0x01, 0xff, 0xc4, 0x00, 0x14, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x08, 0xff, 0xc4, 0x00, 0x14, 0x10, 0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xc4, 0x00,
	0x15, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x07, 0x09, 0xff, 0xc4, 0x00, 0x14, 0x11, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xda, 0x00, 0x0c, 0x03, 0x01,
	0x00, 0x02, 0x11, 0x03, 0x11, 0x00, 0x3f, 0x00, 0x3a, 0x03, 0x15, 0x4d, 0xff, 0xd9,
}

var fetcher = fetch.TestFetcher{
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
	"https://register.eci.ec.europa.eu/core/api/register/logo/8846": image3x1PNG,
	"https://register.eci.ec.europa.eu/core/api/register/logo/8847": image3x1JPG,
	"https://register.eci.ec.europa.eu/core/api/register/details/2024/000009": []byte(`{
		"status": "ONGOING",
		"latestUpdateDate": "24/07/2024 13:52",
		"deadline": "17/05/2025",
		"linguisticVersions": [
			{
				"original": true,
				"languageCode": "EN",
				"title": "Title",
				"objectives": "<p>objectives</p>",
				"annexText": "<ul><li>arg 1: BECAUSE!!!</li><li>arg 2 ...</li></ul>",
				"treaties": "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
				"supportLink": "https://eci.ec.europa.eu/043/public/?lg=fr",
				"website": "https://furfreeeurope.eu/",
				"commissionDecision": {
					"url": "http://eur-lex.europa.eu/legal-content/EN/TXT/PDF/?uri=CELEX:32022D0482&from=EN"
				},
				"additionalDocument": {
					"id": 6729,
					"name": "Fur Free Europe- additional info.pdf",
					"mimeType": "application/pdf",
					"size": 144933
				}
			},
			{
				"original": false,
				"languageCode": "FR",
				"title": "Titre",
				"objectives": "<p>Objectifs</p>",
				"annexText": "<ul><li>arg 1: PARCE QUE!!!</li><li>arg 2 ...</li></ul>",
				"treaties": "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
				"commissionDecision": {
					"document": {
						"id": 8600,
						"name": "CDD-2012000002.pdf",
						"mimeType": "application/pdf",
						"size": 15325
					}
				}
			}
		],
		"categories": [
			{ "categoryType": "SANTE" },
			{ "categoryType": "TRADE" },
			{ "categoryType": "SANTE" },
			{ "categoryType": "AGRI" }
		],
		"progress": [
			{ "name": "REGISTERED", "active": false, "date": "16/03/2022" },
			{ "name": "ANSWERED", "active": true, "date": "07/12/2023" },
			{
				"name": "CLOSED",
				"active": false,
				"date": "01/03/2023",
				"footnoteType": "COLLECTION_EARLY_CLOSURE"
			},
			{ "name": "SUBMITTED", "active": false, "date": "14/06/2023" },
			{ "name": "VERIFICATION", "active": false, "date": "13/03/2023" },
			{
				"name": "COLLECTION_START_DATE",
				"active": false,
				"date": "18/05/2022"
			}
		],
		"sosReport": {
			"updateDate": "24/07/2024",
			"entry": [
				{ "countryCodeType": "FI", "total": 10000 },
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
	wf, to := tool.NewTestTool(fetcher)
	items := fetchIndex(context.Background(), to)
	assert.True(t, wf.NoCalled())
	assert.Equal(t, []indexItem{
		{year: 2024, number: 8},
		{year: 2024, number: 9, logoID: 8846},
	}, items)
}

func TestFetchDetail(t *testing.T) {
	newDate := func(year int, month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, render.DateZone)
	}
	wf, to := tool.NewTestTool(fetcher)
	eci := fetchDetail(context.Background(), to, indexItem{
		year:   2024,
		number: 9,
		logoID: 8846,
	})
	assert.True(t, wf.NoCalled())
	assert.Equal(t, &ECIOut{
		Year:       2024,
		Number:     9,
		LastUpdate: time.Date(2024, time.July, 24, 13, 52, 0, 0, time.UTC),
		Status:     "ONGOING",
		Categorie:  []string{"AGRI", "SANTE", "TRADE"},

		Description: map[language.Language]*Description{
			language.English: {
				Title: "Title",
				SupportLink: &url.URL{
					Scheme:   "https",
					Host:     "eci.ec.europa.eu",
					Path:     "/043/public/",
					RawQuery: "lg=fr",
				},
				Website: &url.URL{
					Scheme: "https",
					Host:   "furfreeeurope.eu",
					Path:   "/",
				},
				PlainDesc: "objectives",
				Objective: "<p>objectives</p>",
				AnnexDoc: &Document{
					URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/6729"),
					Language: language.English,
					Name:     "Fur Free Europe- additional info.pdf",
					Size:     144933,
					MimeType: "application/pdf",
				},
				Annex:  "<ul><li>arg 1: BECAUSE!!!</li><li>arg 2 ...</li></ul>",
				Treaty: "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
			},
			language.French: {
				Title:     "Titre",
				PlainDesc: "Objectifs",
				Objective: "<p>Objectifs</p>",
				AnnexDoc: &Document{
					URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/6729"),
					Language: language.English,
					Name:     "Fur Free Europe- additional info.pdf",
					Size:     144933,
					MimeType: "application/pdf",
				},
				Annex:  "<ul><li>arg 1: PARCE QUE!!!</li><li>arg 2 ...</li></ul>",
				Treaty: "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
			},
		},
		DescriptionOriginalLangage: language.English,

		Timeline: []Timeline{
			{Date: newDate(2022, time.March, 16), Status: "REGISTERED", Register: &[language.Len]*Document{
				language.English: {URL: parseURL("http://eur-lex.europa.eu/legal-content/EN/TXT/PDF/?uri=CELEX:32022D0482&from=EN")},
				language.French: {
					URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/8600"),
					Language: language.French,
					Name:     "CDD-2012000002.pdf",
					MimeType: "application/pdf",
					Size:     15325,
				},
			}},
			{Date: newDate(2022, time.May, 18), Status: "ONGOING"},
			{Date: newDate(2023, time.March, 1), Status: "CLOSED", EarlyClose: true},
			{Date: newDate(2023, time.March, 13), Status: "VERIFICATION"},
			{Date: newDate(2023, time.June, 14), Status: "SUBMITTED"},
			{Date: newDate(2023, time.December, 7), Status: "ANSWERED"},
			{Date: newDate(2025, time.May, 17), Status: "DEADLINE"},
		},

		TotalSignature:        84005,
		PaperSignaturesUpdate: time.Date(2024, time.July, 24, 0, 0, 0, 0, render.DateZone),
		Signature: map[country.Country]uint{
			country.Finland:     10_000,
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
		Threshold:          &threshold_2020_02_01,
		ThresholdRule:      "2020-01-01",
		ThresholdPass:      [country.Len]bool{country.Finland: true},
		ThresholdPassTotal: 1,

		ImageName:   "logo.png",
		ImageWidth:  "3",
		ImageHeight: "1",
		ImageData:   image3x1PNG,
	}, eci)
}

func TestFetchImageJPEG(t *testing.T) {
	wf, to := tool.NewTestTool(fetcher)
	eci := &ECIOut{}
	eci.fetchImage(context.Background(), to, 8847)
	assert.True(t, wf.NoCalled())
	assert.Equal(t, &ECIOut{
		ImageName:   "logo.jpg",
		ImageWidth:  "3",
		ImageHeight: "1",
		ImageData:   image3x1JPG,
	}, eci)
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
