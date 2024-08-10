package service

import (
	"context"
	"sniffle/front"
	"sniffle/service/about"
	"sniffle/service/eu_ec_eci"
	"sniffle/service/home"
	"sniffle/tool"
	"sniffle/tool/render"
)

var List = []tool.Service{
	{Name: "notImplementedPage", Do: notImplementedPage},

	{Name: "about", Do: about.Do},
	{Name: "home", Do: home.Do},
	{Name: "eu_ec_eci", Do: eu_ec_eci.Do},
	{Name: "front", Do: front.Do},
}

func notImplementedPage(_ context.Context, t *tool.Tool) {
	for _, l := range t.Languages {
		t.WriteFile("/eu/index.html", render.Back)
		t.WriteFile("/eu/"+l.String()+".html", render.Back)
		t.WriteFile("/eu/ec/index.html", render.Back)
		t.WriteFile("/eu/ec/"+l.String()+".html", render.Back)
	}
}
