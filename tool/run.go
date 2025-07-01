package tool

import (
	"context"
	"log/slog"
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

	for _, service := range services {
		begin := time.Now()

		if strings.HasPrefix(service.Name, "//") && !DevMode {
			continue
		}

		t := New(config)
		t.Logger = t.Logger.With(slog.Any("service", service.Name))
		service.Do(t)

		t.Log(context.Background(), NoticeLevel, "end", "duration", time.Since(begin))
	}

	end := New(config)

	end.Log(context.Background(), NoticeLevel, "end", "duration", time.Since(globalBegin))
}
