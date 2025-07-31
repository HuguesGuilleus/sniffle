package sch

import (
	"fmt"
	"math"

	"github.com/HuguesGuilleus/sniffle/tool/render"
)

type arrayType struct {
	minLen int
	maxLen int

	item Type
}

// Check the value is an array with any length.
func Array(item Type) Type { return arrayType{0, math.MaxInt, item} }

// Check the value is an array with size length.
func ArraySize(size int, item Type) Type { return arrayType{size, size, item} }

// Check the value is an array with min <= length.
func ArrayMin(min int, item Type) Type { return arrayType{min, math.MaxInt, item} }

// Check the value is an array with min <= length <= max.
func ArrayRange(min, max int, item Type) Type {
	return arrayType{min, max, item}
}

func (t arrayType) Match(v any) error {
	slice, ok := v.([]any)
	if !ok {
		return fmt.Errorf(notArrayFormat, v)
	}

	errs := make(ErrorSlice, 0, len(slice))
	for index, item := range slice {
		if err := t.item.Match(item); err != nil {
			errs.Append(fmt.Sprintf("[%d]", index), err)
		}
	}

	if l := len(slice); l < t.minLen || t.maxLen < l {
		errs = append(errs, fmt.Errorf("expected len in %s but get %d", t.formatRange(), l))
	}

	return errs.Return()
}

func (t arrayType) HTML(indent string) render.Node {
	return render.N("", t.formatRange(), t.item.HTML(indent))
}

func (t arrayType) formatRange() string {
	min := t.minLen
	max := t.maxLen
	if min == max {
		return fmt.Sprintf("[%d]", min)
	} else if min == 0 && max == math.MaxInt {
		return "[]"
	} else if max == math.MaxInt {
		return fmt.Sprintf("[%d..]", min)
	} else {
		return fmt.Sprintf("[%d..%d]", min, max)
	}
}

type emptyArray struct{}

func EmptyArray() Type { return emptyArray{} }

func (emptyArray) Match(v any) error {
	slice, ok := v.([]any)
	if !ok {
		return fmt.Errorf(notArrayFormat, v)
	}

	if len(slice) != 0 {
		return fmt.Errorf("Not empty array, len=%d", len(slice))
	}

	return nil
}

func (emptyArray) HTML(indent string) render.Node {
	return render.N("", "[]")
}
