package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/HuguesGuilleus/sniffle/common"
	"github.com/HuguesGuilleus/sniffle/common/rimage"
	"github.com/HuguesGuilleus/sniffle/front"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/service/about"
	"github.com/HuguesGuilleus/sniffle/service/eu_ec_eci"
	"github.com/HuguesGuilleus/sniffle/service/eu_eca"
	"github.com/HuguesGuilleus/sniffle/service/eu_parl_mep"
	"github.com/HuguesGuilleus/sniffle/service/home"
	"github.com/HuguesGuilleus/sniffle/service/release"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fetch"
	"github.com/HuguesGuilleus/sniffle/tool/render"
	"github.com/HuguesGuilleus/sniffle/tool/writefs"
)

func main() {
	globalBegin := time.Now()

	config := tool.CLI(nil)
	writerSitemap := writefs.Sitemap(&config.Writefile)
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

	writefs.WriteFile(config.Writefile, "/sitemap.txt", writerSitemap.Sitemap(common.Host))

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
