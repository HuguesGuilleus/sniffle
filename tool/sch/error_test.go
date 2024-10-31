package sch

import (
	"bytes"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	buff := bytes.Buffer{}
	logger := slog.New(slog.NewTextHandler(&buff, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				a.Value = slog.TimeValue(time.Time{})
				a.Value = slog.StringValue("")
			}
			return a
		},
	}))

	Log(logger, String("a"), "a")
	Log(logger, String("a"), "b")

	assert.Equal(t, `time="" level=WARN msg=typeCheck err="value \"b\" != \"a\""`+"\n", buff.String())
}

func TestError(t *testing.T) {
	baseErr := errors.New("base")
	we := wrapedError{"b", baseErr}

	assert.Equal(t, baseErr, we.Unwrap())

	errs := make(ErrorSlice, 0)
	errs.Append("a", we)
	assert.Equal(t, "a.b: base", errs.Error())
	assert.EqualValues(t, errs, errs.Unwrap())
	assert.Equal(t, errs, toErrorSlice(errs))
}
