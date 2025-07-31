package tool

import (
	"encoding/json"
	"io"

	"github.com/HuguesGuilleus/sniffle/tool/fetch"
	"github.com/HuguesGuilleus/sniffle/tool/sch"
)

// Fetch and decode in JSON to dto.
//
// Is devmode and ty is not nil, make a first request to check json in type.
func FetchJSON(t *Tool, ty sch.Type, dto any, request *fetch.Request) (fail bool) {
	if DevMode && ty != nil {
		r := t.Fetch(request)
		if r == nil {
			return true
		}
		defer r.Body.Close()

		var value any
		dec := json.NewDecoder(r.Body)
		dec.UseNumber()
		dec.Decode(&value)
		sch.Log(t.Logger.With("url", request.URL), ty, value)
	}

	response := t.Fetch(request)
	if response == nil {
		return true
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(dto); err != nil {
		t.Warn("http.decodeJsonFail", "id", request.ID(), "url", request.URL.String(), "err", err.Error())
		return true
	}

	return false
}

func FetchAll(t *Tool, request *fetch.Request) []byte {
	response := t.Fetch(request)
	if response == nil {
		return nil
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Warn("http.readAllfail", "url", request.URL.String(), "err", err.Error())
		return nil
	}

	return data
}
