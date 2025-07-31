package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sniffle/common"
	"sniffle/common/rimage"
	"sniffle/front"
	"sniffle/front/translate"
	"sniffle/service/about"
	"sniffle/service/eu_ec_eci"
	"sniffle/service/eu_eca"
	"sniffle/service/eu_parl_mep"
	"sniffle/service/home"
	"sniffle/service/release"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/render"
	"sniffle/tool/writefile"
	"time"
)

func main() {
	globalBegin := time.Now()

	config := tool.CLI(nil)
	writerSitemap := writefile.Sitemap(&config.Writefile)
	config.LongTasksMap[rimage.NameResizeJpeg] = rimage.FetchResizeJpeg

	// Service call
	config.Run("front", front.Do)
	config.Run("notImplementedPage", notImplementedPage)

	config.Run("about", about.Do)
	config.Run("release", release.Do)
	config.Run("home", home.Do)

	config.Run("eu_ec_eci", eu_ec_eci.Do)
	config.Run("//eu_eca", eu_eca.Do)
	config.Run("//eu_parl_mep", eu_parl_mep.Do)

	config.Writefile.WriteFile("/sitemap.txt", writerSitemap.Sitemap(common.Host))

	// Make cache debug index.
	if tool.DevMode {
		err := fetch.Debug(tool.CLICache(), func(m *fetch.Meta) int {
			switch m.URL.Host {
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

	slog.New(config.LogHandler).Log(context.Background(), tool.NoticeLevel, "end", "duration", time.Since(globalBegin))
}

func notImplementedPage(t *tool.Tool) {
	t.WriteFile("/eu/index.html", render.Back)
	t.WriteFile("/eu/ec/index.html", render.Back)
	t.WriteFile("/eu/parl/index.html", render.Back)
	for _, l := range translate.Langs {
		t.WriteFile(l.Path("/eu/"), render.Back)
		t.WriteFile(l.Path("/eu/ec/"), render.Back)
		t.WriteFile(l.Path("/eu/parl/"), render.Back)
	}
}
