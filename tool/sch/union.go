package sch

import "sniffle/tool/render"

/* AND */

type andType []Type

func And(types ...Type) Type { return andType(types) }

func (and andType) Match(v any) error {
	errs := make(ErrorSlice, 0, len(and))
	for _, t := range and {
		err := t.Match(v)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs.Return()
}

func (and andType) HTML(indent string) render.Node {
	return render.N("", render.S(and, " & ", func(t Type) render.Node {
		return t.HTML(indent)
	}))
}

/* OR */

type orType []Type

// Check the value checks one of the Type.
func Or(types ...Type) Type { return orType(types) }

func (or orType) Match(v any) error {
	errs := make(ErrorSlice, 0, len(or))
	for _, t := range or {
		err := t.Match(v)
		if err == nil {
			return nil
		}
		errs = append(errs, err)
	}
	return errs
}

func (or orType) HTML(indent string) render.Node {
	return render.N("", render.S(or, " | ", func(t Type) render.Node {
		return t.HTML(indent)
	}))
}
