// eu_ec_eci service for European Citizens' Initiative.
package eu_ec_eci

import (
	"fmt"
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
)

func Do(t *tool.Tool) {
	eciByYear := fetchAllAcepted(t)
	t.WriteFile("/eu/ec/eci/index.html", lredirect.All)
	t.WriteFile("/eu/ec/eci/data/index.html", render.Back)
	t.WriteFile("/eu/ec/eci/data/extradelay/index.html", lredirect.All)
	t.WriteFile("/eu/ec/eci/data/threshold/index.html", lredirect.All)
	t.WriteFile("/eu/ec/eci/schema.html", schemaPage)
	for _, l := range translate.Langs {
		writeIndex(t, eciByYear, l)
		t.WriteFile(l.Path("/eu/ec/eci/data/"), render.Back)
		renderDataExtraDelay(t, l)
		renderDataThreshold(t, l)
	}
	for year, eciSlice := range eciByYear {
		t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/index.html", year), render.Back)
		for _, eci := range eciSlice {
			redirect := lredirect.Page(fmt.Sprintf("https://citizens-initiative.europa.eu/initiatives/details/%d/%06d", eci.Year, eci.Number), eci.Langs())
			t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/index.html", eci.Year, eci.Number), redirect)
			for _, l := range translate.Langs {
				if eci.Description[l] != nil {
					writeOne(t, eci, l)
				} else {
					t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/%s.html", eci.Year, eci.Number, l), redirect)
				}
			}
			if img := eci.Image; img != nil {
				t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/logo%s", eci.Year, eci.Number, img.Raw.Extension), img.Raw.Data)
				if res := img.Resized; res != nil {
					t.WriteFile(fmt.Sprintf("/eu/ec/eci/%d/%d/logo%s", eci.Year, eci.Number, res.Extension), res.Data)
				}
			}
		}
	}

	refused := fetchRefusedAll(t)
	for _, eci := range refused {
		redirect := lredirect.Page(eci.OfficielLink(), eci.Langs())
		p := "/eu/ec/eci/refused/" + printUint(eci.ID) + "/"
		t.WriteFile(p+"index.html", redirect)
		for _, l := range translate.Langs {
			if eci.Lang == l {
				t.WriteFile(l.Path(p), renderRefusedOne(eci))
			} else {
				t.WriteFile(l.Path(p), redirect)
			}
		}
	}
	t.WriteFile("/eu/ec/eci/refused/index.html", lredirect.All)
	for _, l := range translate.Langs {
		t.WriteFile(l.Path("/eu/ec/eci/refused/"), renderRefusedIndex(refused, t.HostURL+"/eu/ec/eci/refused/", l))
	}
}
