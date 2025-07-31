package eu_parl_mep

import (
	"encoding/xml"
	"fmt"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"time"
)

const (
	metaURL = "https://data.europarl.europa.eu/OdpDatasetService/Datasets/members-of-the-european-parliament-meps-parliamentary-term%d"
	dataURL = "https://data.europarl.europa.eu/distribution/meps_%d_%d.rdf"
)

type mep struct {
	Identifier uint
	Label      string

	Term []int
}

func fetchVersion(t *tool.Tool, term int) (version int) {
	dtoMeta := struct {
		OdpDatasetVersions []struct {
			VersionLabel int `json:",string"`
		}
	}{}
	if tool.FetchJSON(t, metaType, &dtoMeta, fetch.Fmt(metaURL, term)) {
		return -1
	}

	for _, op := range dtoMeta.OdpDatasetVersions {
		version = max(version, op.VersionLabel)
	}

	return
}

func fetchMep(t *tool.Tool, term int) (list []mep) {
	u := fmt.Sprintf(dataURL, term, fetchVersion(t, term))
	data := tool.FetchAll(t, fetch.URL(u))

	dto := struct {
		Person []struct {
			Identifier uint   `xml:"identifier"`
			Label      string `xml:"label"`
			GivenName  string `xml:"givenName"`
			FamilyName string `xml:"familyName"`

			Bday bday `xml:"bday"`
			// placeOfBirth
			// 	<person:citizenship rdf:resource="http://publications.europa.eu/resource/authority/country/FRA"/>

			// 	<vcard:hasEmail rdf:resource="mailto:manon.aubry@europarl.europa.eu"/>
		}
	}{}
	err := xml.Unmarshal(data, &dto)
	if err != nil {
		t.Error("xml.decode", "url", u, "err", err.Error())
	}

	list = make([]mep, len(dto.Person))
	for i, dto := range dto.Person {
		list[i] = mep{
			Identifier: dto.Identifier,
			Label:      dto.Label,
			Term:       []int{term},
		}
	}
	return
}

type bday struct {
	Time time.Time
}

func (p *bday) UnmarshalText(text []byte) error {
	t, err := time.Parse(time.DateOnly, string(text))
	if err != nil {
		return err
	}
	p.Time = t
	return nil
}
