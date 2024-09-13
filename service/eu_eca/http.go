package eu_eca

import (
	"context"
	"net/http"
	"sniffle/tool"
)

func Do(ctx context.Context, t *tool.Tool) {
	if !tool.DevMode {
		t.Info("skip")
		return
	}

	if tool.DevMode {
		t.WriteFile("/eu/eca/src-token.json", tool.FetchAll(ctx, t, http.MethodPost, "https://www.eca.europa.eu/_api/contextinfo", http.Header{
			"Accept":       []string{"application/json; odata=verbose"},
			"Content-Type": []string{"application/json; odata=verbose"},
		}, nil))
	}
	v := struct {
		D struct {
			GetContextWebInformation struct {
				FormDigestValue string
			}
		}
	}{}
	tool.FetchJSON(ctx, t, http.MethodPost, "https://www.eca.europa.eu/_api/contextinfo", http.Header{
		"Accept":       []string{"application/json; odata=verbose"},
		"Content-Type": []string{"application/json; odata=verbose"},
	}, nil, &v)
	token := v.D.GetContextWebInformation.FormDigestValue

	if tool.DevMode {
		t.WriteFile("/eu/eca/src-doc.json", tool.FetchAll(ctx, t, http.MethodPost, "https://www.eca.europa.eu/_vti_bin/ECA.Internet/DocSetService.svc/SearchDocs", http.Header{
			"Accept":          []string{"application/json"},
			"Content-Type":    []string{"application/json"},
			"X-RequestDigest": []string{token},
		}, []byte(`{"searchInput":{"RowLimit":1000,"Filter":"ECADocType=\"Annual report\"","Lang":"French","LangCode":"FR"}}`)))
	}
}
