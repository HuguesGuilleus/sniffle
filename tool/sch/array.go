package sch

import (
	"fmt"
	"sniffle/tool/render"
)

type arrayType struct {
	item Type
}

func Array(item Type) Type { return arrayType{item} }

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

	return errs.Return()
}

func (t arrayType) HTML(indent string) render.Node {
	return render.N("", `[]`, t.item.HTML(indent))
}
