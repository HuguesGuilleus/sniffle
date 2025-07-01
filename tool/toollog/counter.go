package toollog

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"
)

type failCounterHandler struct {
	err  *uint64
	warn *uint64
	sub  slog.Handler
}

func (h *failCounterHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.sub.Enabled(ctx, l)
}

func (h *failCounterHandler) Handle(ctx context.Context, r slog.Record) error {
	switch {
	case slog.LevelError <= r.Level:
		atomic.AddUint64(h.err, 1)
	case slog.LevelWarn <= r.Level:
		atomic.AddUint64(h.warn, 1)
	}
	return h.sub.Handle(ctx, r)
}

func (h *failCounterHandler) WithAttrs(attr []slog.Attr) slog.Handler {
	return &failCounterHandler{
		err:  h.err,
		warn: h.warn,
		sub:  h.sub.WithAttrs(attr),
	}
}

func (h *failCounterHandler) WithGroup(group string) slog.Handler {
	return &failCounterHandler{
		err:  h.err,
		warn: h.warn,
		sub:  h.sub.WithGroup(group),
	}
}

// Replace the slog handler with a wraper that count warn and error.
// The returned function return a short string with time and full sucess or fail count.
func CountFail(sub *slog.Handler) func() []byte {
	h := &failCounterHandler{
		err:  new(uint64),
		warn: new(uint64),
		sub:  *sub,
	}
	*sub = h

	return func() []byte {
		now := time.Now().UTC()
		if *h.warn+*h.err == 0 {
			return now.AppendFormat(nil, "2006-01-02 15:04:05 | FULL SUCESS\r\n")
		} else {
			return fmt.Appendf(nil, "%s | WARN(%d), ERROR(%d)\r\n", now.Format(time.DateTime), *h.warn, *h.err)
		}
	}
}
