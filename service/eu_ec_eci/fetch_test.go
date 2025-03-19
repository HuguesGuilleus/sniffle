package eu_ec_eci

import (
	"net/url"
	"sniffle/common"
	"sniffle/common/country"
	"sniffle/common/language"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/fetch"
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

var fetcher = map[string]*fetch.TestResponse{
	acceptedIndexURL: fetch.TRs(200, `{
		"entries": [
			{
				"id": 8846,
				"year": "2024",
				"number": "000008"
			},
			{
				"id": 8845,
				"year": "2024",
				"number": "000009"
			}
		]
	}`),
	"https://register.eci.ec.europa.eu/core/api/register/logo/8846": fetch.TR(200, image3x1PNG),
	"https://register.eci.ec.europa.eu/core/api/register/logo/8847": fetch.TR(200, image3x1JPG),
	"https://register.eci.ec.europa.eu/core/api/register/details/2024/000009": fetch.TRs(200, `{
		"status": "ONGOING",
		"latestUpdateDate": "24/07/2024 13:52",
		"deadline": "17/05/2025",
		"logo": {"id": 8846},
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
					"celex": "32022D0482",
					"corrigendum": "32020D0674R(01)"
				},
				"additionalDocument": {
					"id": 6729,
					"name": "Fur Free Europe- additional info.pdf",
					"mimeType": "application/pdf",
					"size": 144933
				},
				"draftLegal": {
					"id": 8936,
					"name": "DRAFT LEGAL ACT.pdf",
					"mimeType": "application/pdf",
					"size": 47930
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
		"members": [
				{
					"type": "LEGAL_ENTITY",
					"fullName": "Zavod za zaščito in napredek reproduktivnih pravic My Voice, My Choice ",
					"email": "https://www.myvoice-mychoice.org",
					"residenceCountry": "si",
					"privacyApplied": false
				},
				{
					"type": "REPRESENTATIVE",
					"fullName": "Remo NANNETTI",
					"email": "",
					"residenceCountry": "it",
					"privacyApplied": false,
					"replacedMember": [
						{
							"type": "REPRESENTATIVE",
							"fullName": "Gabriele BONCI",
							"email": "stopcrueltystopslaughter@gmail.com",
							"residenceCountry": "it",
							"privacyApplied": true,
							"startingDate": "24/07/2024",
							"endDate": "03/09/2024"
						}
					],
					"startingDate": "03/09/2024"
				}
		],
		"categories": [
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
		"funding": {
			"lastUpdate": "15/05/2023",
			"sponsors": [
				{
					"name": "Campaigns and Activism for Animals in the Industry",
					"date": "01/03/2023",
					"amount": 11986,
					"privateSponsor": false,
					"anonymized": false
				},
				{
					"name": "[ANONYMIZED]",
					"date": "19/08/2024",
					"amount": 881.1,
					"privateSponsor": true,
					"anonymized": true
				}
			],
			"document": {
				"id": 9122,
				"name": "2023_03_01_ffe_financial reporting.pdf",
				"mimeType": "application/pdf",
				"size": 500553
			},
			"totalAmount": 12867.1
		},
		"sosReport": {
			"updateDate": "24/07/2024",
			"totalSignatures": 88783,
			"entry": [
				{ "countryCodeType": "FI", "total": 10000, "afterSubmission": true },
				{ "countryCodeType": "RO", "total": 1362 },
				{ "countryCodeType": "CY", "total": 60 },
				{ "countryCodeType": "LU", "total": 354 },
				{ "countryCodeType": "PL", "total": 1879 },
				{ "countryCodeType": "IT", "total": 4940 },
				{ "countryCodeType": "HR", "total": 628 },
				{ "countryCodeType": "SE", "total": 16420 },
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
		},
		"answer": {
			"links": [
				{
					"defaultLanguageCode": "EN",
					"defaultName": "COMMUNICATION",
					"defaultLink": "https://citizens-initiative.europa.eu/sites/default/files/2023-12/C_2023_8362_EN.pdf ",
					"link": [
						{
							"languageCode": "FR",
							"link": "https://eur-lex.europa.eu/legal-content/FR/TXT/?uri=CELEX:52023XC01559"
						},
						{
							"languageCode": "EN",
							"link": "https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:52023XC01559"
						}
					]
				},
				{
					"defaultLanguageCode": "EN",
					"defaultName": "ANNEX",
					"defaultLink": "https://ec.europa.eu/transparency/documents-register/detail?ref=C(2023)4489&lang=en",
					"link": [
						{
							"languageCode": "EN",
							"link": "https://ec.europa.eu/transparency/documents-register/detail?ref=C(2023)4489&lang=en"
						},
						{
							"languageCode": "FR",
							"link": "https://ec.europa.eu/transparency/documents-register/detail?ref=C(2023)4489&lang=fr"
						}
					]
				},
				{
					"defaultLanguageCode": "EN",
					"defaultName": "PRESS_RELEASE",
					"defaultLink": "https://ec.europa.eu/commission/presscorner/detail/en/ip_23_6251"
				},
				{
					"defaultLanguageCode": "EN",
					"defaultName": "FOLLOW_UP",
					"defaultLink": "https://citizens-initiative.europa.eu/fur-free-europe"
				}
			]
		}
	}`),
}

func TestFetchIndex(t *testing.T) {
	wf, to := tool.NewTestTool(fetcher)
	items := fetchAcceptedIndex(to)
	assert.True(t, wf.NoCalled())
	assert.Equal(t, []indexItem{
		{id: 8846, year: 2024, number: 8},
		{id: 8845, year: 2024, number: 9},
	}, items)
}

func TestFetchDetail(t *testing.T) {
	defer func(langs []language.Language) { translate.Langs = langs }(translate.Langs)
	translate.Langs = []language.Language{language.English, language.French}

	wf, to := tool.NewTestTool(fetcher)
	eci := fetchDetail(to, indexItem{
		id:     648,
		year:   2024,
		number: 9,
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
				Website: &url.URL{
					Scheme: "https",
					Host:   "furfreeeurope.eu",
					Path:   "/",
				},
				SupportLink: parseURL("https://eci.ec.europa.eu/043/public/?lg=fr"),
				FollowUp:    parseURL("https://citizens-initiative.europa.eu/fur-free-europe_en"),
				PlainDesc:   "objectives",
				Objective:   "<p>objectives</p>",
				AnnexDoc: &Document{
					URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/6729"),
					Language: language.English,
					Name:     "Fur Free Europe- additional info.pdf",
					Size:     144933,
					MimeType: "application/pdf",
				},
				DraftLegal: &Document{
					URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/8936"),
					Language: language.English,
					Name:     "DRAFT LEGAL ACT.pdf",
					Size:     47930,
					MimeType: "application/pdf",
				},
				Annex:  "<ul><li>arg 1: BECAUSE!!!</li><li>arg 2 ...</li></ul>",
				Treaty: "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
			},
			language.French: {
				Title:     "Titre",
				FollowUp:  parseURL("https://citizens-initiative.europa.eu/fur-free-europe_fr"),
				PlainDesc: "Objectifs",
				Objective: "<p>Objectifs</p>",
				AnnexDoc: &Document{
					URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/6729"),
					Language: language.English,
					Name:     "Fur Free Europe- additional info.pdf",
					Size:     144933,
					MimeType: "application/pdf",
				},
				DraftLegal: &Document{
					URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/8936"),
					Language: language.English,
					Name:     "DRAFT LEGAL ACT.pdf",
					Size:     47930,
					MimeType: "application/pdf",
				},
				Annex:  "<ul><li>arg 1: PARCE QUE!!!</li><li>arg 2 ...</li></ul>",
				Treaty: "Articolo 6 lettera a), Articolo 114, Articolo 168, Articolo 169 TFUE",
			},
		},
		OriginalLangage: language.English,

		Timeline: []Event{
			{
				Date: newDate(2022, time.March, 16), Status: "REGISTERED",
				Register: &[language.Len]*Document{
					language.English: {URL: parseURL("https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32022D0482")},
					language.French: {
						URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/8600"),
						Language: language.French,
						Name:     "CDD-2012000002.pdf",
						MimeType: "application/pdf",
						Size:     15325,
					},
				},
				RegisterCorrigendum: &[language.Len]*Document{
					language.English: {URL: parseURL("https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX%3A32020D0674R%2801%29")},
				},
			},
			{Date: newDate(2022, time.May, 18), Status: "COLLECTION_START_DATE"},
			{Date: newDate(2023, time.March, 1), Status: "CLOSED", EarlyClose: true, ExtraDelay: []ExtraDelay{extraDelay_2021_1121, extraDelay_2021_3879}},
			{Date: newDate(2023, time.March, 13), Status: "VERIFICATION"},
			{Date: newDate(2023, time.June, 14), Status: "SUBMITTED"},
			{
				Date: newDate(2023, time.December, 7), Status: "ANSWERED",
				AnswerPressRelease: docs(&Document{
					Language: language.English,
					URL:      parseURL("https://ec.europa.eu/commission/presscorner/detail/en/ip_23_6251"),
				}),
				AnswerResponse: docs(&Document{
					Language: language.English,
					URL:      parseURL("https://citizens-initiative.europa.eu/sites/default/files/2023-12/C_2023_8362_EN.pdf"),
				}, &Document{
					Language: language.French,
					URL:      parseURL("https://eur-lex.europa.eu/legal-content/FR/TXT/?uri=CELEX:52023XC01559"),
				}, &Document{
					Language: language.English,
					URL:      parseURL("https://eur-lex.europa.eu/legal-content/EN/TXT/?uri=CELEX:52023XC01559"),
				}),
				AnswerAnnex: docs(&Document{
					Language: language.English,
					URL:      parseURL("https://ec.europa.eu/transparency/documents-register/detail?ref=C(2023)4489&lang=en"),
				}, &Document{
					Language: language.French,
					URL:      parseURL("https://ec.europa.eu/transparency/documents-register/detail?ref=C(2023)4489&lang=fr"),
				}),
			},
			{Date: newDate(2025, time.May, 17), Status: "DEADLINE"},
		},

		Signature: []Signature{
			{country.Austria, 1099, false, 13_395, false},
			{country.Belgium, 2101, false, 14_805, false},
			{country.Bulgaria, 196, false, 11_985, false},
			{country.Croatia, 628, false, 8_460, false},
			{country.Cyprus, 60, false, 4_230, false},
			{country.Czechia, 367, false, 14_805, false},
			{country.Denmark, 866, false, 9_870, false},
			{country.Estonia, 244, false, 4_935, false},
			{country.Finland, 10_000, true, 9_870, false},
			{country.France, 24384, false, 55_695, false},
			{country.Germany, 12985, false, 67_680, false},
			{country.Greece, 551, false, 14_805, false},
			{country.Hungary, 630, false, 14_805, false},
			{country.Ireland, 1757, false, 9_165, false},
			{country.Italy, 4940, false, 53_580, false},
			{country.Latvia, 149, false, 5_640, false},
			{country.Lithuania, 427, false, 7_755, false},
			{country.Luxembourg, 354, false, 4_230, false},
			{country.Malta, 61, false, 4_230, false},
			{country.Netherlands, 3856, false, 20_445, false},
			{country.Poland, 1879, false, 36_660, false},
			{country.Portugal, 650, false, 14_805, false},
			{country.Romania, 1362, false, 23_265, false},
			{country.Slovakia, 434, false, 9_870, false},
			{country.Slovenia, 775, false, 5_640, false},
			{country.Spain, 11608, false, 41_595, false},
			{country.Sweden, 16420, false, 14_805, true},
		},
		TotalSignature:        88783,
		PaperSignaturesUpdate: time.Date(2024, time.July, 24, 0, 0, 0, 0, render.DateZone),
		ThresholdRule:         "2020-01-01",
		ThresholdPassTotal:    1,

		Image: &common.ResizedImage{
			Raw: common.Image{
				Extension: ".png",
				Width:     "3",
				Height:    "1",
				Data:      image3x1PNG,
			},
		},

		// Members
		Members: []Member{
			{
				Type:             "LEGAL_ENTITY",
				FullName:         "Zavod za zaščito in napredek reproduktivnih pravic My Voice, My Choice",
				DisplayURL:       "https://www.myvoice-mychoice.org",
				HrefURL:          "https://www.myvoice-mychoice.org",
				ResidenceCountry: country.Slovenia,
			},
			{
				Type:             "REPRESENTATIVE",
				FullName:         "Remo NANNETTI",
				ResidenceCountry: country.Italy,
				Start:            newDate(2024, time.September, 3),
				Replaced: &Member{
					Type:             "REPRESENTATIVE",
					FullName:         "Gabriele BONCI",
					HrefURL:          "mailto:stopcrueltystopslaughter@gmail.com",
					DisplayURL:       "stopcrueltystopslaughter@gmail.com",
					ResidenceCountry: country.Italy,
					Start:            newDate(2024, time.July, 24),
					End:              newDate(2024, time.September, 3),
					Privacy:          true,
				},
			},
		},

		// Funding
		FundingUpdate: time.Date(2023, time.May, 15, 0, 0, 0, 0, render.DateZone),
		FundingTotal:  12867.1,
		FundingDocument: &Document{
			URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/9122"),
			Name:     "2023_03_01_ffe_financial reporting.pdf",
			Size:     500553,
			MimeType: "application/pdf",
		},
		Sponsor: []Sponsor{
			{
				Name:      "Campaigns and Activism for Animals in the Industry",
				IsPrivate: false,
				Amount:    11986,
				Date:      time.Date(2023, time.March, 1, 0, 0, 0, 0, render.DateZone),
			},
			{
				Name:      "",
				IsPrivate: true,
				Amount:    881.1,
				Date:      time.Date(2024, time.August, 19, 0, 0, 0, 0, render.DateZone),
			},
		},
	}, eci)
}

func TestFetchImageJPEG(t *testing.T) {
	wf, to := tool.NewTestTool(fetcher)
	assert.Equal(t, &common.ResizedImage{
		Raw: common.Image{
			Extension: ".jpg",
			Width:     "3",
			Height:    "1",
			Data:      image3x1JPG,
		},
	}, fetchImage(to, 8847))
	assert.Nil(t, fetchImage(to, 0))
	assert.True(t, wf.NoCalled())
}

func newDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, render.DateZone)
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func docs(base *Document, specific ...*Document) *[language.Len]*Document {
	slice := new([language.Len]*Document)
	for i := range slice {
		slice[i] = base
	}
	for _, doc := range specific {
		slice[doc.Language] = doc
	}
	return slice
}
