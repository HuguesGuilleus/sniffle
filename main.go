package main

import (
	"flag"
	"sniffle/common/resize0"
	"sniffle/front"
	"sniffle/front/translate"
	"sniffle/service/about"
	"sniffle/service/eu_ec_eci"
	"sniffle/service/eu_eca"
	"sniffle/service/home"
	"sniffle/service/release"
	"sniffle/tool"
	"sniffle/tool/render"
	"time"
)

func main() {
	host := flag.String("host", "https://sniffle.eu/", "The host absolute URL")

	config := tool.CLI(map[string]time.Duration{"": time.Millisecond * 100})

	config.HostURL = *host
	config.Languages = translate.Langs

	config.LongTasksMap = map[string]func([]byte) ([]byte, error){
		resize0.Name: resize0.Resize,
	}

	tool.Run(config,
		tool.Service{Name: "front", Do: front.Do},
		tool.Service{Name: "notImplementedPage", Do: notImplementedPage},

		tool.Service{Name: "about", Do: about.Do},
		tool.Service{Name: "release", Do: release.Do},
		tool.Service{Name: "home", Do: home.Do},

		tool.Service{Name: "eu_ec_eci", Do: eu_ec_eci.Do},
		tool.Service{Name: "//eu_eca", Do: eu_eca.Do},
	)
}

func notImplementedPage(t *tool.Tool) {
	t.WriteFile("/eu/index.html", render.Back)
	t.WriteFile("/eu/ec/index.html", render.Back)
	for _, l := range translate.Langs {
		t.WriteFile(l.Path("/eu/"), render.Back)
		t.WriteFile(l.Path("/eu/ec/"), render.Back)
	}
}
