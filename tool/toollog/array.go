package toollog

import (
	"context"
	"log/slog"
)

type slice []slog.Handler

// A handlers slices that call each sub handlers.
func Slice(handlers ...slog.Handler) *slog.Logger {
	return slog.New(slice(handlers))
}

func (s slice) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range s {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (s slice) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range s {
		if !h.Enabled(ctx, r.Level) {
			continue
		} else if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (s slice) WithAttrs(attrs []slog.Attr) slog.Handler {
	n := make(slice, len(s))
	for i, h := range s {
		n[i] = h.WithAttrs(attrs)
	}
	return n
}

func (s slice) WithGroup(name string) slog.Handler {
	n := make(slice, len(s))
	for i, h := range s {
		n[i] = h.WithGroup(name)
	}
	return n
}
