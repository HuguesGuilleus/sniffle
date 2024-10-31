package sch

import (
	"errors"
	"log/slog"
)

// Check if the value is math the Type t.
// If not, log all errors with message:"typeCheck"
func Log(logger *slog.Logger, t Type, value any) {
	for _, err := range toErrorSlice(t.Match(value)) {
		logger.Warn("typeCheck", "err", err.Error())
	}
}

type wrapedError struct {
	key string
	err error
}

func (w wrapedError) Error() string { return w.key + ": " + w.err.Error() }
func (w wrapedError) Unwrap() error { return w.err }

type ErrorSlice []error

func (s ErrorSlice) Error() string   { return errors.Join(s...).Error() }
func (s ErrorSlice) Unwrap() []error { return s }

func (s *ErrorSlice) Append(key string, err error) {
	for _, err := range toErrorSlice(err) {
		werr, ok := err.(wrapedError)
		if !ok {
			werr = wrapedError{key, err}
		} else {
			werr.key = key + "." + werr.key
		}
		*s = append(*s, werr)
	}
}

// Return nil if the slice do not contain error.
func (s ErrorSlice) Return() error {
	if len(s) != 0 {
		return s
	}
	return nil
}

func toErrorSlice(err error) ErrorSlice {
	if err == nil {
		return nil
	}
	errs, ok := err.(ErrorSlice)
	if !ok {
		return ErrorSlice{err}
	}
	return errs
}
