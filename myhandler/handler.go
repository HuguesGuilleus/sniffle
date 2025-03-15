// myhandler is a slog text handler.
package myhandler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"unicode"
)

const timeFormat = "2006-01-02T15:04:05"

type handler struct {
	level slog.Level

	w io.Writer
	m *sync.Mutex

	service string
	attr    bytes.Buffer
	group   bytes.Buffer
}

func New(w io.Writer, level slog.Level) slog.Handler {
	return &handler{
		level: level,
		w:     w,
		m:     new(sync.Mutex),
	}
}

func (h *handler) Enabled(_ context.Context, l slog.Level) bool { return l >= h.level }

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	n := h.clone()
	g := h.group.String()
	for _, attr := range attrs {
		if attr.Key == "service" && attr.Value.Kind() == slog.KindString && n.service == "" {
			n.service = attr.Value.String()
		} else {
			printAttr(&n.attr, attr, g)
		}
	}
	return n
}

func (h *handler) WithGroup(name string) slog.Handler {
	n := h.clone()
	quote(&n.group, name)
	n.group.WriteByte('.')
	return n
}

func (h *handler) clone() (n *handler) {
	n = new(handler)
	n.level = h.level
	n.w = h.w
	n.m = h.m
	n.service = h.service
	if h.attr.Len() != 0 {
		n.attr.Write(h.attr.Bytes())
	}
	if h.group.Len() != 0 {
		n.group.Write(h.group.Bytes())
	}
	return
}

func (h *handler) Handle(_ context.Context, record slog.Record) error {
	buff := bytes.Buffer{}
	buff.WriteString(record.Time.UTC().Format("15:04:05 "))
	buff.WriteString(record.Level.String())
	for range 16 - buff.Len() {
		buff.WriteByte('_')
	}
	buff.WriteByte(' ')
	buff.WriteByte('[')
	if h.service != "" {
		quote(&buff, h.service)
		buff.WriteByte('|')
	}
	quote(&buff, record.Message)
	buff.WriteByte(']')

	if h.level < slog.LevelInfo {
		if f := runtime.FuncForPC(record.PC); f != nil {
			file, line := f.FileLine(record.PC)
			fmt.Fprintf(&buff, " src=%s:%d", filepath.Base(file), line)
		}
	}

	buff.Write(h.attr.Bytes())

	g := h.group.String()
	record.Attrs(func(attr slog.Attr) bool {
		printAttr(&buff, attr, g)
		return true
	})
	buff.WriteByte('\n')

	h.m.Lock()
	defer h.m.Unlock()
	buff.WriteTo(h.w)

	return nil
}

func printAttr(buff *bytes.Buffer, attr slog.Attr, group string) {
	buff.WriteByte(' ')
	buff.WriteString(group)
	quote(buff, attr.Key)
	buff.WriteByte('=')

	switch v := attr.Value; attr.Value.Kind() {
	case slog.KindBool:
		if v.Bool() {
			buff.WriteString("true_")
		} else {
			buff.WriteString("false")
		}

	case slog.KindDuration:
		buff.WriteString(v.Duration().String())
	case slog.KindTime:
		buff.WriteString(v.Time().UTC().Format(timeFormat))

	case slog.KindFloat64, slog.KindUint64, slog.KindInt64:
		buff.WriteString(v.String())

	case slog.KindGroup:
		buff.WriteByte('[')
		for _, a := range v.Group() {
			printAttr(buff, a, group)
		}
		buff.WriteByte(' ')
		buff.WriteByte(']')

	default:
		quote(buff, v.String())
	}
}

func quote(buff *bytes.Buffer, s string) {
	for _, r := range s {
		if r == '\'' || r == '\\' || r == '"' || unicode.IsSpace(r) {
			buff.WriteString(strconv.Quote(s))
			return
		}
	}
	buff.WriteString(s)
}
