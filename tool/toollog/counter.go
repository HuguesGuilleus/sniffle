package toollog

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"
)

type FailCounterHandler struct {
	warn uint64
	err  uint64
}

func (counter *FailCounterHandler) Enabled(_ context.Context, l slog.Level) bool {
	return slog.LevelWarn <= l
}

func (counter *FailCounterHandler) Handle(_ context.Context, r slog.Record) error {
	l := r.Level
	switch {
	case slog.LevelError <= l:
		atomic.AddUint64(&counter.err, 1)
	case slog.LevelWarn <= l:
		atomic.AddUint64(&counter.warn, 1)
	}
	return nil
}

func (counter *FailCounterHandler) WithAttrs([]slog.Attr) slog.Handler { return counter }
func (counter *FailCounterHandler) WithGroup(string) slog.Handler      { return counter }

// A basic format with time and sucess or number of failure.
func (counter *FailCounterHandler) Bytes() []byte {
	now := time.Now().UTC()
	if counter.warn+counter.err == 0 {
		return now.AppendFormat(nil, "2006-01-02 15:04:05 | FULL SUCESS\r\n")
	} else {
		return fmt.Appendf(nil, "%s | WARN(%d), ERROR(%d) \r\n", now.Format(time.DateTime), counter.warn, counter.err)
	}
}
