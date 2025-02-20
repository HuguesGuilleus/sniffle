package tool

import (
	"context"
	"log/slog"
	"sniffle/tool/toollog"
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

	defer func(oldLogger *slog.Logger) { config.Logger = oldLogger }(config.Logger)
	counterLogHandler := toollog.FailCounterHandler{}
	config.Logger = toollog.Slice(config.Logger.Handler(), &counterLogHandler)

	htmlFiles := make([]string, 0)
	writeSum := uint64(0)

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

		writeSum += t.writeSum
	}

	end := New(config)
	end.writeSitemap(htmlFiles)
	end.WriteFile("/log", counterLogHandler.Bytes())
	writeSum += end.writeSum

	config.Logger.Log(context.Background(), NoticeLevel, "end", "duration", time.Since(globalBegin), "writeSum", writeSum)
}

func (t *Tool) writeSitemap(paths []string) {
	data := make([]byte, 0)

	for _, p := range paths {
		data = append(data, t.HostURL...)
		data = append(data, p...)
		data = append(data, '\n')
	}

	t.WriteFile("/sitemap.txt", data)
}
