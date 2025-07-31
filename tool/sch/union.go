package sch

import (
	"strings"

	"github.com/HuguesGuilleus/sniffle/tool/render"
)

/* AND */

type andType []Type

func And(types ...Type) Type {
	if len(types) <= 1 {
		panic("sch.And() must pass multiples types")
	}
	return andType(types)
}

func (and andType) Match(v any) error {
	errs := make(ErrorSlice, 0, len(and))
	for _, t := range and {
		err := t.Match(v)
		if err != nil {
			errs.Append("$and", err)
		}
	}
	return errs.Return()
}

func (and andType) HTML(indent string) render.Node {
	return render.N("", render.S(and, " & ", func(t Type) render.Node {
		return t.HTML(indent)
	}))
}

/* ENUM string */

type enumString struct {
	orType
	s string
}

// A enum of string.
// Equal Or(String(...), ...)
func EnumString(enums ...string) TypeStringer {
	or := make(orType, len(enums))
	for i, e := range enums {
		or[i] = String(e)
	}
	return enumString{or, strings.Join(enums, "|")}
}

func (enum enumString) String() string { return enum.s }

/* OR */

type orType []Type

// Check the value checks one of the Type.
func Or(types ...Type) Type {
	if len(types) <= 1 {
		panic("sch.Or() must pass multiples types")
	}
	return orType(types)
}

func (or orType) Match(v any) error {
	errs := make(ErrorSlice, 0, len(or))
	for _, t := range or {
		err := t.Match(v)
		if err == nil {
			return nil
		}
		errs.Append("$or", err)
	}
	return errs
}

func (or orType) HTML(indent string) render.Node {
	return render.N("", render.S(or, " | ", func(t Type) render.Node {
		return t.HTML(indent)
	}))
}
