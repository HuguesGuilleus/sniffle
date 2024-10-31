package sch

import (
	"fmt"
	"regexp"
	"sniffle/tool/render"
	"time"
)

/* REGEXP */

type regexpType struct {
	*regexp.Regexp
}

func Regexp(str string) TypeStringer { return regexpType{regexp.MustCompile(str)} }

func (r regexpType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	if !r.Regexp.MatchString(s) {
		return fmt.Errorf("rexpexp /%s/ not match %q", r.Regexp, s)
	}

	return nil
}

func (r regexpType) HTML(_ string) render.Node {
	return render.N(xstringMarkup, "regexp/", render.N("u", r.String()), "/")
}

/* TIME */

type timeType struct {
	layout string
}

// Check that the value is a string formated with layout.
func Time(layout string) Type { return timeType{layout} }

func (t timeType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	_, err := time.Parse(t.layout, s)
	return err
}

func (t timeType) HTML(_ string) render.Node {
	return render.Na(xstringMarkup, "title", "A time value encoded into a string").N(
		"string(", render.N("u", t.layout), ")",
	)
}
