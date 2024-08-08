package myhandler

import (
	"bytes"
	"context"
	"log/slog"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMyHandler(t *testing.T) {
	ptr := reflect.ValueOf(staticPC).Pointer()
	now := time.Date(2024, time.August, 3, 16, 29, 24, 354, time.UTC)
	buff := &bytes.Buffer{}
	ctx := context.Background()

	assert.False(t, New(buff, slog.LevelWarn).Enabled(ctx, slog.LevelInfo))
	assert.True(t, New(buff, slog.LevelWarn).Enabled(ctx, slog.LevelWarn))
	assert.True(t, New(buff, slog.LevelWarn).Enabled(ctx, slog.LevelError))

	buff.Reset()
	r0 := slog.NewRecord(now, slog.LevelWarn, "yolo", ptr)
	New(buff, slog.LevelWarn).Handle(ctx, r0)
	assert.Equal(t, "2024-08-03T16:29:24 WARN [yolo]\n", buff.String())
	buff.Reset()
	New(buff, slog.LevelDebug).Handle(ctx, r0)
	assert.Equal(t, "2024-08-03T16:29:24 WARN [yolo] src=f_test.go:3\n", buff.String())

	r1 := slog.NewRecord(now, slog.LevelWarn, "y\"ol\"o", ptr)
	r1.Add("int", -53)
	buff.Reset()
	New(buff, slog.LevelWarn).
		WithGroup("g1").
		WithAttrs([]slog.Attr{slog.Any(" with", true)}).
		WithGroup("g2").
		Handle(ctx, r1)
	assert.Equal(t, "2024-08-03T16:29:24 WARN [\"y\\\"ol\\\"o\"] g1.\" with\"=true_ g1.g2.int=-53\n", buff.String())

	r2 := slog.NewRecord(now, slog.LevelInfo+1, "m", ptr)
	r2.Add("u", uint64(1234))
	r2.Add("f64", 4.1)
	r2.Add("false", false)
	r2.Add("nil", nil)
	r2.Add("dur", time.Second*2)
	r2.Add("now", now)
	r2.Add("gr", slog.GroupValue(slog.Any("k1", 1), slog.Any("k2", 2)))
	buff.Reset()
	New(buff, slog.LevelInfo).Handle(ctx, r2)
	assert.Equal(t, "2024-08-03T16:29:24 INFO+1 [m] u=1234 f64=4.1 false=false nil=<nil> dur=2s now=2024-08-03T16:29:24 gr=[ k1=1 k2=2 ]\n", buff.String())
}
