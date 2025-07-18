package eu_parl_meps

import (
	"encoding/xml"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"time"
)

type meps struct {
	Identifier uint
	Label      string

	Term []uint
}

func fetchMeps(t *tool.Tool) (list []meps) {
	// f, _ := os.Create("service/eu_parl_meps/a.xml")
	// defer f.Close()
	// io.Copy(f, t.Fetch(fetch.URL("https://data.europarl.europa.eu/distribution/meps_10_45.rdf")).Body)
	// return

	// t.Fetch(fetch.URL("https://data.europarl.europa.eu/OdpDatasetService/Datasets/members-of-the-european-parliament-meps-parliamentary-term10"))

	u := "https://data.europarl.europa.eu/distribution/meps_10_45.rdf"
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

	list = make([]meps, len(dto.Person))
	for i, dto := range dto.Person {
		list[i] = meps{
			Identifier: dto.Identifier,
			Label:      dto.Label,
			Term:       []uint{10},
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
