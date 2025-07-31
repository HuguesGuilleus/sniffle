package eu_ec_eci

import (
	"sniffle/common/language"
	"sniffle/tool"
	"sniffle/tool/render"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchRefused(t *testing.T) {
	wfs, to := tool.NewTestTool(fetcher)
	defer assert.Empty(t, wfs)
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
