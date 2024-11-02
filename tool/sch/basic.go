package sch

import (
	"encoding/json"
	"fmt"
	"math"
	"sniffle/tool/render"
	"strconv"
	"strings"
)

const (
	notArrayFormat      = "not map: %+v"
	notBoolFormat       = "not bool: %+v"
	notJsonNumberFormat = "not json number: %+v"
	notMapFormat        = "not map: %+v"
	notNilFormat        = "not nil: %+v"
	notStringFormat     = "not string: %+v"

	baseMarkup    = "span.sch-base"
	intMarkup     = "span.sch-int"
	stringMarkup  = "span.sch-str"
	xstringMarkup = "span.sch-xstr"
)

/* ANY */

type anyType struct{}

// Every value is valid.
func Any() Type { return anyType{} }

func (anyType) Match(any) error { return nil }

func (anyType) HTML(_ string) render.Node { return render.N(baseMarkup, "...") }

/* ANY BOOL */

type anyBoolType struct{}

// Check is the value is a bool.
func AnyBool() Type { return anyBoolType{} }

func (anyBoolType) Match(v any) error {
	_, ok := v.(bool)
	if !ok {
		return fmt.Errorf(notBoolFormat, v)
	}
	return nil
}

func (anyBoolType) HTML(_ string) render.Node { return render.N(baseMarkup, "boolean") }

/* BOOL */

type boolType bool

// Check if the value is true.
func True() Type { return boolType(true) }

// Check if the value is false.
func False() Type { return boolType(false) }

func (t boolType) Match(v any) error {
	b, ok := v.(bool)
	if !ok {
		return fmt.Errorf(notBoolFormat, v)
	} else if bool(t) != b {
		return fmt.Errorf("expected bool %t != received %t", bool(t), b)
	}
	return nil
}

func (t boolType) HTML(_ string) render.Node {
	if t {
		return render.N(baseMarkup, "true")
	} else {
		return render.N(baseMarkup, "false")
	}
}

/* As INT */

type asIntType struct {
	min, max int64
}

// Check the value is a json.Number integer.
func AsAnyInt() Type { return asIntType{math.MinInt64, math.MaxInt64} }

// Check the value v is a json.Number integer and 0 <= v.
func AsPositiveInt() Type { return asIntType{0, math.MaxInt64} }

// Check the value v is a json.Number integer and 0 < v.
func AsStrictPositiveInt() Type { return asIntType{1, math.MaxInt64} }

// Check the value v is a json.Number integer and v <= 0.
func AsNegativeInt() Type { return asIntType{math.MinInt64, 0} }

// Check the value v is a json.Number integer and v < 0.
func AsStrictNegativeInt() Type { return asIntType{math.MinInt64, -1} }

// Check the value v is a json.Number integer and min <= v <= max.
func IntervalAsInt(min, max int64) Type { return asIntType{min, max} }

func (t asIntType) Match(v any) error {
	j, ok := v.(json.Number)
	if !ok {
		return fmt.Errorf(notJsonNumberFormat, v)
	}
	i, err := j.Int64()
	if err != nil {
		return fmt.Errorf("json.Number integer %q: %w", j, err)
	}
	if i < t.min {
		return fmt.Errorf("value %d < min %d", i, t.min)
	} else if t.max < i {
		return fmt.Errorf("max %d < value %d", t.max, i)
	}
	return nil
}

func (t asIntType) HTML(_ string) render.Node {
	return render.Na(intMarkup, "title", "Integer").N(t.htmlContent())
}

func (t asIntType) htmlContent() string {
	switch {
	case t.min == math.MinInt64 && t.max == math.MaxInt64:
		return "integer"
	case t.min == math.MinInt64:
		return fmt.Sprintf("[ .. %d ]", t.max)
	case t.max == math.MaxInt64:
		return fmt.Sprintf("[ %d .. ]", t.min)
	default:
		return fmt.Sprintf("[ %d .. %d ]", t.min, t.max)
	}
}

/* AS STRING INT */

type asStringIntType struct {
	asIntType
}

// Check the value is a integer into a string.
func AsAnyStringInt() TypeStringer { return IntervalAsStringInt(math.MinInt64, math.MaxInt64) }

// Check the value v is a integer into a string and 0 <= v.
func AsPositiveStringInt() TypeStringer { return IntervalAsStringInt(0, math.MaxInt64) }

// Check the value v is a integer into a string and 0 < v.
func AsStrictPositiveStringInt() TypeStringer { return IntervalAsStringInt(1, math.MaxInt64) }

// Check the value v is a integer into a string and v <= 0.
func AsNegativeStringInt() TypeStringer { return IntervalAsStringInt(math.MinInt64, 0) }

// Check the value v is a integer into a string and v < 0.
func AsStrictNegativeStringInt() TypeStringer { return IntervalAsStringInt(math.MinInt64, -1) }

// Check the value v is a integer into a string and min <= v <= max.
func IntervalAsStringInt(min, max int64) TypeStringer {
	return asStringIntType{asIntType{min, max}}
}

func (t asStringIntType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	return t.asIntType.Match(json.Number(s))
}

func (t asStringIntType) HTML(_ string) render.Node {
	return render.Na(intMarkup, "title", "Integer into a string").N("string(", t.htmlContent(), ")")
}

func (t asStringIntType) String() string {
	return "string(" + t.htmlContent() + ")"
}

/* Nil */

type nilType struct{}

func Nil() Type { return nilType{} }

func (nilType) Match(v any) error {
	if v != nil {
		return fmt.Errorf(notNilFormat, v)
	}
	return nil
}

func (nilType) HTML(_ string) render.Node { return render.N(baseMarkup, "null") }

/* ANY STRING */

type anyStringType struct{}

// The value is a string.
func AnyString() Type { return anyStringType{} }

func (anyStringType) Match(v any) error {
	_, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	return nil
}

func (anyStringType) HTML(_ string) render.Node { return render.N(baseMarkup, "string") }

/* NOT NIL STRING */

type notEmptyStringType struct{}

func NotEmptyString() Type { return notEmptyStringType{} }

func (notEmptyStringType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("value %+v is not a empty string", v)
	} else if s == "" {
		return fmt.Errorf("value %q is not a empty string", s)
	}
	return nil
}

func (notEmptyStringType) HTML(_ string) render.Node { return render.N(baseMarkup, "not-empty-string") }

/* ARROUND STRING */

type arroundStringType string

func ArroundString(s string) TypeStringer {
	s = strings.ToLower(strings.TrimSpace(s))
	return arroundStringType(s)
}

func (t arroundStringType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	s = strings.ToLower(strings.TrimSpace(s))
	if s != string(t) {
		return fmt.Errorf("value %q != %q", s, t)
	}
	return nil
}

func (t arroundStringType) HTML(_ string) render.Node {
	return render.N(stringMarkup, "~", strconv.Quote(string(t)))
}

func (t arroundStringType) String() string { return "~" + string(t) }

/* VALUE STRING */

type stringType string

func String(s string) TypeStringer { return stringType(s) }

func (t stringType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	} else if s != string(t) {
		return fmt.Errorf("value %q != %q", s, t)
	}
	return nil
}

func (t stringType) HTML(_ string) render.Node {
	return render.N(stringMarkup, strconv.Quote(string(t)))
}

func (t stringType) String() string { return string(t) }

/* PRINT */

type printlnType string

// Print in the console the match argument with the id.
// Use this type only for development.
func Print(id string) Type { return printlnType(id) }

func (id printlnType) Match(v any) error {
	fmt.Printf("[%s] %+v\n", id, v)
	return nil
}

func (id printlnType) HTML(_ string) render.Node { return render.N(baseMarkup, "[log:", id, "]") }
