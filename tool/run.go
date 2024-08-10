package tool

import (
	"context"
	"log/slog"
	"time"
)

type Service struct {
	Name string
	Do   func(context.Context, *Tool)
}

func Run(ctx context.Context, config *Config, services []Service) {
	htmlFiles := make([]string, 0)

	for _, service := range services {
		begin := time.Now()

		t := New(config)
		t.Logger = t.Logger.With(slog.Any("S", service.Name))
		t.htmlFiles = htmlFiles
		service.Do(ctx, t)

		t.Info("duration", "duration", time.Since(begin))
		htmlFiles = t.htmlFiles
	}

	New(config).writeSitemap(htmlFiles)
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
