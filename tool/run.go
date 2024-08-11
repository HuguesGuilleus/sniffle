package tool

import (
	"context"
	"log/slog"
	"time"
)

const NoticeLevel = slog.LevelInfo + 2

type Service struct {
	Name string
	Do   func(context.Context, *Tool)
}

func Run(ctx context.Context, config *Config, services []Service) {
	globalBegin := time.Now()

	htmlFiles := make([]string, 0)

	writeSum := uint64(0)

	for _, service := range services {
		begin := time.Now()

		t := New(config)
		t.Logger = t.Logger.With(slog.Any("service", service.Name))
		t.htmlFiles = htmlFiles
		service.Do(ctx, t)

		t.Log(ctx, NoticeLevel, "end", "duration", time.Since(begin))
		htmlFiles = t.htmlFiles

		writeSum += t.writeSum
	}

	to := New(config)
	to.writeSitemap(htmlFiles)
	writeSum += to.writeSum

	config.Logger.Log(ctx, NoticeLevel, "end", "duration", time.Since(globalBegin), "writeSum", writeSum)
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
