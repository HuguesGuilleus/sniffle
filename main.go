package main

import (
	"context"
	"flag"
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
	"sniffle/service/home"
	"sniffle/service/release"
	"sniffle/tool"
	"sniffle/tool/fetch"
	"sniffle/tool/render"
	"sniffle/tool/toollog"
	"sniffle/tool/writefile"
	"time"
)

func main() {
	globalBegin := time.Now()

	config := tool.CLI(nil)
	writerSitemap := writefile.Sitemap(&config.Writefile)
	config.LongTasksMap[rimage.NameResizeJpeg] = rimage.FetchResizeJpeg

	defer func(f func() []byte) {
		config.Writefile.WriteFile("/log", f())
	}(toollog.CountFail(&config.LogHandler))

	// Service call
	config.Run("front", front.Do)
	config.Run("notImplementedPage", notImplementedPage)

	config.Run("about", about.Do)
	config.Run("release", release.Do)
	config.Run("home", home.Do)

	config.Run("eu_ec_eci", eu_ec_eci.Do)
	config.Run("//eu_eca", eu_eca.Do)

	config.Writefile.WriteFile("/sitemap.txt", writerSitemap.Sitemap(common.Host))

	// Make cache debug index.
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

	slog.New(config.LogHandler).Log(context.Background(), tool.NoticeLevel, "end", "duration", time.Since(globalBegin))
}

func notImplementedPage(t *tool.Tool) {
	t.WriteFile("/eu/index.html", render.Back)
	t.WriteFile("/eu/ec/index.html", render.Back)
	for _, l := range translate.Langs {
		t.WriteFile(l.Path("/eu/"), render.Back)
		t.WriteFile(l.Path("/eu/ec/"), render.Back)
	}
}
