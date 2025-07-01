package tool

import (
	"context"
	"log/slog"
	"sniffle/common"
	"strings"
	"time"
)

const NoticeLevel = slog.LevelInfo + 2

type Service struct {
	// If name begin by "//" and we are not in DevMode, skip it.
	Name string
	Do   func(*Tool)
}

func Run(config *Config, services ...Service) {
	globalBegin := time.Now()

	htmlFiles := make([]string, 0)

	for _, service := range services {
		begin := time.Now()

		if strings.HasPrefix(service.Name, "//") && !DevMode {
			continue
		}

		t := New(config)
		t.Logger = t.Logger.With(slog.Any("service", service.Name))
		t.htmlFiles = htmlFiles
		service.Do(t)

		t.Log(context.Background(), NoticeLevel, "end", "duration", time.Since(begin))
		htmlFiles = t.htmlFiles
	}

	end := New(config)
	end.writeSitemap(htmlFiles)

	end.Log(context.Background(), NoticeLevel, "end", "duration", time.Since(globalBegin))
}

func (t *Tool) writeSitemap(paths []string) {
	data := make([]byte, 0)

	for _, p := range paths {
		data = append(data, []byte(common.Host)...)
		data = append(data, p...)
		data = append(data, '\n')
	}

	t.WriteFile("/sitemap.txt", data)
}
