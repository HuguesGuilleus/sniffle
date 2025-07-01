package main

import (
	"flag"
	"fmt"
	"os"
	"sniffle/common/rimage"
	"sniffle/front"
	"sniffle/front/translate"
	"sniffle/service/about"
	"sniffle/service/eu_ec_eci"
	"sniffle/service/eu_eca"
	"sniffle/service/home"
	"sniffle/service/release"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/render"
	"sniffle/tool/toollog"
)

func main() {
	config := tool.CLI(nil)
	config.LongTasksMap[rimage.NameResizeJpeg] = rimage.FetchResizeJpeg

	defer func(f func() []byte) {
		config.Writefile.WriteFile("/log", f())
	}(toollog.CountFail(&config.LogHandler))

	tool.Run(config,
		tool.Service{Name: "front", Do: front.Do},
		tool.Service{Name: "notImplementedPage", Do: notImplementedPage},

		tool.Service{Name: "about", Do: about.Do},
		tool.Service{Name: "release", Do: release.Do},
		tool.Service{Name: "home", Do: home.Do},

		tool.Service{Name: "eu_ec_eci", Do: eu_ec_eci.Do},
		tool.Service{Name: "//eu_eca", Do: eu_eca.Do},
	)

	if tool.DevMode {
		err := fetch.Debug(flag.CommandLine.Lookup("cache").Value.String(), func(host string) int {
			switch host {
			case "register.eci.ec.europa.eu":
				return fetch.DebugKeepIndex
			default:
				return fetch.DebugKeepData
			}
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func notImplementedPage(t *tool.Tool) {
	t.WriteFile("/eu/index.html", render.Back)
	t.WriteFile("/eu/ec/index.html", render.Back)
	for _, l := range translate.Langs {
		t.WriteFile(l.Path("/eu/"), render.Back)
		t.WriteFile(l.Path("/eu/ec/"), render.Back)
	}
}
