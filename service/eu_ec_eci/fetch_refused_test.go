package eu_ec_eci

import (
	"sniffle/common/language"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/render"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchRefused(t *testing.T) {
	fetcher := map[string]*fetch.TestResponse{
		"https://register.eci.ec.europa.eu/core/api/register/search/REFUSED/EN/0/0": fetch.TRs(200, `{
			"entries": [
				{"id": 42}
			]
		}`),
		"https://register.eci.ec.europa.eu/core/api/register/details/42": fetch.TRs(200, `{
			"refusalDate": "30/04/2019",
			"refusalDocument": {
				"id": 4373,
				"name": "COM_2019_3305_public.pdf",
				"mimeType": "application/pdf",
				"size": 230455
			},
			"linguisticVersions": [{
				"languageCode": "EN",
				"title": "Stopping trade with Israeli settlements operating in the Occupied Palestinian Territory",
				"objectives": "<p>Hello",
				"treaties": "Treaty on the Functioning of the European Union, 2012, at Article 207.",
				"website": "https://github.com/",
				"additionalDocument": {
					"id": 4313,
					"name": "add_doc.pdf",
					"mimeType": "application/pdf",
					"size": 456845
				},
				"draftLegal": {
					"id": 4314,
					"name": "draft_legal.pdf",
					"mimeType": "application/pdf",
					"size": 180804
				},
				"commissionDecision": {
					"celex": "32019D0722"
				}
			}]
		}`),
	}

	wf, to := tool.NewTestTool(fetcher)
	defer assert.True(t, wf.NoCalled())
	assert.Equal(t, []*ECIRefused{
		{
			ID:         42,
			Lang:       language.English,
			Website:    parseURL("https://github.com/"),
			Title:      "Stopping trade with Israeli settlements operating in the Occupied Palestinian Territory",
			PlainDesc:  "Hello",
			Objectives: "<p>Hello</p>",
			Treaties:   "Treaty on the Functioning of the European Union, 2012, at Article 207.",

			AnnexDoc: &Document{
				URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/4313"),
				Language: language.English,
				Name:     "add_doc.pdf",
				Size:     456845,
				MimeType: "application/pdf",
			},
			DraftLegal: &Document{
				URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/4314"),
				Language: language.English,
				Name:     "draft_legal.pdf",
				Size:     180804,
				MimeType: "application/pdf",
			},

			RefusedDate: time.Date(2019, time.April, 30, 0, 0, 0, 0, render.DateZone),
			RefusalDocument: Document{
				URL:      parseURL("https://register.eci.ec.europa.eu/core/api/register/document/4373"),
				Language: language.English,
				Name:     "COM_2019_3305_public.pdf",
				Size:     230455,
				MimeType: "application/pdf",
			},
			RefusedCELEX: "32019D0722",
		},
	}, fetchRefusedAll(to))
}
