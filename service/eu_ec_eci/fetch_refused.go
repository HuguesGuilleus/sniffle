package eu_ec_eci

import (
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fetch"
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/HuguesGuilleus/sniffle/tool/securehtml"
)

func fetchRefusedAll(t *tool.Tool) []*ECIRefused {
	index := fetchRefusedIndex(t)

	refused := make([]*ECIRefused, 0, len(index))
	wg := sync.WaitGroup{}
	wg.Add(len(index))
	mutex := sync.Mutex{}
	for _, entry := range index {
		go func(id uint) {
			defer wg.Done()
			eci := fetchOneRefused(t, id)
			if eci == nil {
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			refused = append(refused, eci)
		}(entry.id)
	}
	wg.Wait()

	slices.SortFunc(refused, func(a, b *ECIRefused) int {
		return b.RefusedDate.Compare(a.RefusedDate)
	})

	return refused
}

type ECIRefused struct {
	ID   uint
	Lang language.Language

	Website    *url.URL
	Title      string
	PlainDesc  string
	Objectives render.H
	AnnexText  render.H
	Treaties   string

	AnnexDoc   *Document
	DraftLegal *Document

	RefusedDate     time.Time
	RefusalDocument Document
	RefusedCELEX    string
}

func fetchOneRefused(t *tool.Tool, id uint) *ECIRefused {
	idString := strconv.FormatUint(uint64(id), 10)
	request := fetch.URL("https://register.eci.ec.europa.eu/core/api/register/details/" + idString)
	if tool.DevMode {
		t.WriteFile("/eu/ec/eci/refused/"+idString+"/src.json", tool.FetchAll(t, request))
	}
	dto := struct {
		RefusalDocument docDTO  `json:"refusalDocument"`
		RefusedDate     dtoDate `json:"refusalDate"`
		Desc            [1]struct {
			Lang       language.Language `json:"languageCode"`
			WebSite    string            `json:"website"`
			Title      string            `json:"title"`
			Objectives string            `json:"objectives"`
			AnnexText  string            `json:"annexText"`
			Treaties   string            `json:"treaties"`

			AnnexDoc   *docDTO `json:"additionalDocument"`
			DraftLegal *docDTO `json:"draftLegal"`

			CommissionDecision struct {
				CELEX string `json:"celex"`
			} `json:"commissionDecision"`
		} `json:"linguisticVersions"`
	}{}
	if tool.FetchJSON(t, refusedOneType, &dto, request) {
		return nil
	}
	desc := dto.Desc[0]

	if desc.AnnexDoc != nil && desc.DraftLegal != nil {
		if desc.AnnexDoc.Name == desc.DraftLegal.Name &&
			desc.AnnexDoc.Size == desc.DraftLegal.Size &&
			desc.AnnexDoc.MimeType == desc.DraftLegal.MimeType {
			desc.DraftLegal = nil
		}
	}

	return &ECIRefused{
		ID:   id,
		Lang: desc.Lang,

		Website:    securehtml.ParseURL(desc.WebSite),
		Title:      desc.Title,
		PlainDesc:  securehtml.Text(desc.Objectives, plainDescMax),
		Objectives: securehtml.Secure(desc.Objectives),
		AnnexText:  securehtml.Secure(desc.AnnexText),
		Treaties:   desc.Treaties,

		AnnexDoc:   desc.AnnexDoc.Document(desc.Lang),
		DraftLegal: desc.DraftLegal.Document(desc.Lang),

		RefusedDate:     dto.RefusedDate.Time.In(render.DateZone),
		RefusalDocument: *dto.RefusalDocument.Document(desc.Lang),
		RefusedCELEX:    desc.CommissionDecision.CELEX,
	}
}

func (eci *ECIRefused) OfficielLink() string {
	return "https://citizens-initiative.europa.eu/initiatives/details/" + printUint(eci.ID) + "_" + eci.Lang.String()
}

func (eci *ECIRefused) Langs() []language.Language {
	for _, l := range translate.Langs {
		if eci.Lang == l {
			return []language.Language{l}
		}
	}
	return nil
}
