package eu_eca

import (
	"net/http"
	"sniffle/tool"
	"sniffle/tool/fetch"
)

func Do(t *tool.Tool) {
	if !tool.DevMode {
		t.Info("skip")
		return
	}

	if tool.DevMode {
		t.WriteFile("/eu/eca/src-token.json", tool.FetchAll(t,
			fetch.R(http.MethodPost, "https://www.eca.europa.eu/_api/contextinfo", nil,
				"Accept", "application/json; odata=verbose",
				"Content-Type", "application/json; odata=verbose",
			),
		))
	}
	v := struct {
		D struct {
			GetContextWebInformation struct {
				FormDigestValue string
			}
		}
	}{}
	tool.FetchJSON(t, nil, &v, fetch.R(http.MethodPost, "https://www.eca.europa.eu/_api/contextinfo", nil,
		"Accept", "application/json; odata=verbose",
		"Content-Type", "application/json; odata=verbose",
	))
	token := v.D.GetContextWebInformation.FormDigestValue

	if tool.DevMode {
		t.WriteFile("/eu/eca/src-doc.json", tool.FetchAll(t, fetch.Rs(http.MethodPost, "https://www.eca.europa.eu/_vti_bin/ECA.Internet/DocSetService.svc/SearchDocs", `{"searchInput":{"RowLimit":1000,"Filter":"ECADocType=\"Annual report\"","Lang":"French","LangCode":"FR"}}`,
			"Accept", "application/json",
			"Content-Type", "application/json",
			"X-RequestDigest", token,
		)))
	}
}
