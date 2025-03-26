// eu_ec_eci service for European Citizens' Initiative.
package eu_ec_eci

import (
	"cmp"
	"fmt"
	"slices"
	"sniffle/front/lredirect"
	"sniffle/front/translate"
	"sniffle/tool"
	"sniffle/tool/render"
	"sync"
)

func Do(t *tool.Tool) {
	eciByYear := fetchAllAcepted(t)
	t.WriteFile("/eu/ec/eci/index.html", lredirect.All)
	t.WriteFile("/eu/ec/eci/data/index.html", render.Back)
	t.WriteFile("/eu/ec/eci/data/extradelay/index.html", lredirect.All)
	t.WriteFile("/eu/ec/eci/data/threshold/index.html", lredirect.All)
	t.WriteFile("/eu/ec/eci/schema.html", schemaPage)
	for _, l := range translate.Langs {
		renderIndex(t, eciByYear, l)
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
					renderOne(t, eci, l)
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

	ExploreRefused(t)
}

func fetchAllAcepted(t *tool.Tool) map[uint][]*ECIOut {
	infoIndex := fetchAcceptedIndex(t)
	wg := sync.WaitGroup{}
	wg.Add(1 + len(infoIndex))
	go func() {
		defer wg.Done()
		checkThreashold(t)
	}()
	eciByYear := make(map[uint][]*ECIOut)
	mutex := sync.Mutex{}
	for _, info := range infoIndex {
		go func(info indexItem) {
			defer wg.Done()
			eci := fetchDetail(t, info)
			if eci == nil {
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			eciByYear[eci.Year] = append(eciByYear[eci.Year], eci)
		}(info)
	}
	wg.Wait()

	for _, byYear := range eciByYear {
		slices.SortFunc(byYear, func(a, b *ECIOut) int {
			return cmp.Compare(b.Number, a.Number)
		})
	}

	return eciByYear
}
