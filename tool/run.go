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
	for _, service := range services {
		begin := time.Now()

		t := New(config)
		t.Logger = t.Logger.With(slog.Any("S", service.Name))
		service.Do(ctx, t)

		t.Info("duration", "duration", time.Since(begin))
	}
}
