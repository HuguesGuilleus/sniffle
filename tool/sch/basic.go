package sch

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/HuguesGuilleus/sniffle/tool/render"
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
	floatMarkup   = "span.sch-float"
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

/* AS FLOAT64 */

type floatType struct {
	min, max float64
}

// Check if the value is a json.Number float64.
func AnyFloat() Type { return floatType{-math.MaxFloat64, math.MaxFloat64} }

// Check the value f is a json.Number float64 and 0 <= f.
func PositiveFloat() Type { return floatType{0, math.MaxFloat64} }

// Check the value f is a json.Number float64 and 0 < f.
func StrictPositiveFloat() Type {
	return floatType{math.SmallestNonzeroFloat64, math.MaxFloat64}
}

// Check the value f is a json.Number float64 and f <= 0.
func NegativeFloat() Type { return floatType{-math.MaxFloat64, 0} }

// Check the value f is a json.Number float64 and f < 0.
func StrictNegativeFloat() Type {
	return floatType{-math.MaxFloat64, -math.SmallestNonzeroFloat64}
}

// Check the value f is a json.Number float64 and min <= f <= max
func IntervalFloat(min, max float64) Type { return floatType{min, max} }

// Check the value f is json.Number float64 and equal to c.
func ConstFloat(c float64) Type { return floatType{c, c} }

func (t floatType) Match(v any) error {
	j, ok := v.(json.Number)
	if !ok {
		return fmt.Errorf(notJsonNumberFormat, v)
	}

	f, err := j.Float64()
	if err != nil {
		return fmt.Errorf("json.Number float64 %q: %w", j, err)
	}

	if f < t.min {
		return fmt.Errorf("value %f < min %f", f, t.min)
	} else if t.max < f {
		return fmt.Errorf("max %f < value %f", t.max, f)
	}

	return nil
}

func (t floatType) HTML(_ string) render.Node {
	return render.Na(floatMarkup, "title", "Float64").N(t.htmlContent())
}
func (t floatType) htmlContent() render.H {
	format := func(f float64) string {
		if f == 0 {
			return "0.0"
		}
		if math.Trunc(f) == f {
			return strconv.FormatFloat(f, 'f', 1, 64)
		}
		return strconv.FormatFloat(f, 'f', -1, 64)
	}
	switch {
	case t.min == -math.MaxFloat64 && t.max == math.MaxFloat64:
		return "float64"
	case t.min == -math.MaxFloat64 && t.max == -math.SmallestNonzeroFloat64:
		return render.H(fmt.Sprintf("( .. < 0.0 )"))
	case t.min == math.SmallestNonzeroFloat64 && t.max == math.MaxFloat64:
		return render.H(fmt.Sprintf("( 0.0 > .. )"))
	case t.min == -math.MaxFloat64:
		return render.H(fmt.Sprintf("( .. %s )", format(t.max)))
	case t.max == math.MaxFloat64:
		return render.H(fmt.Sprintf("( %s .. )", format(t.min)))
	case t.max == t.min:
		return render.H(format(t.min))
	default:
		return render.H(fmt.Sprintf("( %s .. %s )", format(t.min), format(t.max)))
	}
}

/* AS STRING FLOAT */

type stringFloatType struct {
	floatType
}

// Check the value is a float64 into a string.
func AnyStringFloat() TypeStringer {
	return stringFloatType{
		floatType: floatType{-math.MaxFloat64, math.MaxFloat64},
	}
}

// Check the value is a float64 >= 0.0 into a string.
func PositiveStringFloat() TypeStringer {
	return stringFloatType{
		floatType: floatType{0, math.MaxFloat64},
	}
}

// Check the value is a float64 > 0.0 into a string.
func StrictPositiveStringFloat() TypeStringer {
	return stringFloatType{
		floatType: floatType{math.SmallestNonzeroFloat64, math.MaxFloat64},
	}
}

// Check the value is a float64 <= 0.0 into a string.
func NegativeStringFloat() TypeStringer {
	return stringFloatType{
		floatType: floatType{-math.MaxFloat64, 0},
	}
}

// Check the value is a float64 < 0.0 into a string.
func StrictNegativeStringFloat() TypeStringer {
	return stringFloatType{
		floatType: floatType{-math.MaxFloat64, -math.SmallestNonzeroFloat64},
	}
}

// Check the value is a float64 f with min <= f <= max into a string.
func IntervalStringFloat(min, max float64) TypeStringer {
	return stringFloatType{floatType: floatType{min, max}}
}

// Check the value is a float64 equal c into a string.
func ConstStringFloat(c float64) TypeStringer {
	return stringFloatType{floatType: floatType{c, c}}
}

func (t stringFloatType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	return t.floatType.Match(json.Number(s))
}

func (t stringFloatType) HTML(_ string) render.Node {
	return render.Na(floatMarkup, "title", "Float64 into a string").N("string(", t.htmlContent(), ")")
}

func (t stringFloatType) String() string {
	return "string(" + string(t.floatType.htmlContent()) + ")"
}

/* INT */

type intType struct {
	min, max int64
}

// Check the value is a json.Number integer.
func AnyInt() Type { return intType{math.MinInt64, math.MaxInt64} }

// Check the value v is a json.Number integer and 0 <= v.
func PositiveInt() Type { return intType{0, math.MaxInt64} }

// Check the value v is a json.Number integer and 0 < v.
func StrictPositiveInt() Type { return intType{1, math.MaxInt64} }

// Check the value v is a json.Number integer and v <= 0.
func NegativeInt() Type { return intType{math.MinInt64, 0} }

// Check the value v is a json.Number integer and v < 0.
func StrictNegativeInt() Type { return intType{math.MinInt64, -1} }

// Check the value v is a json.Number integer and min <= v <= max.
func IntervalInt(min, max int64) Type { return intType{min, max} }

// Check the value i is a json.Number integer and equal to c.
func ConstInt(c int64) Type { return intType{c, c} }

func (t intType) Match(v any) error {
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

func (t intType) HTML(_ string) render.Node {
	return render.Na(intMarkup, "title", "Integer").N(t.htmlContent())
}

func (t intType) htmlContent() string {
	switch {
	case t.min == math.MinInt64 && t.max == math.MaxInt64:
		return "integer"
	case t.min == t.max:
		return fmt.Sprintf("%d", t.min)
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
	intType
}

// Check the value is a integer into a string.
func AnyStringInt() TypeStringer { return IntervalStringInt(math.MinInt64, math.MaxInt64) }

// Check the value v is a integer into a string and 0 <= v.
func PositiveStringInt() TypeStringer { return IntervalStringInt(0, math.MaxInt64) }

// Check the value v is a integer into a string and 0 < v.
func StrictPositiveStringInt() TypeStringer { return IntervalStringInt(1, math.MaxInt64) }

// Check the value v is a integer into a string and v <= 0.
func NegativeStringInt() TypeStringer { return IntervalStringInt(math.MinInt64, 0) }

// Check the value v is a integer into a string and v < 0.
func StrictNegativeStringInt() TypeStringer { return IntervalStringInt(math.MinInt64, -1) }

// Check the value v is a integer into a string and min <= v <= max.
func IntervalStringInt(min, max int64) TypeStringer {
	return asStringIntType{intType{min, max}}
}

func (t asStringIntType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	return t.intType.Match(json.Number(s))
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

func (id printlnType) HTML(_ string) render.Node { return render.N(baseMarkup, "[print:", id, "]") }
