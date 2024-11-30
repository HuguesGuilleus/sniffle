package sch

import (
	"errors"
	"fmt"
	"sniffle/tool/render"
)

type mapType struct {
	fields      []MapField
	extraFields bool
}
type MapField struct {
	fieldKey   TypeStringer
	fieldValue Type
	required   bool
	comments   []string
}

func Map(fields ...MapField) Type      { return &mapType{fields, false} }
func MapExtra(fields ...MapField) Type { return &mapType{fields, true} }

// Add a field to the map.
func Field(key TypeStringer, value Type, required bool) MapField {
	return MapField{key, value, required, nil}
}

// Add a required field to the map.
// The key is a String() Type.
func FieldSR(key string, value Type) MapField {
	return Field(String(key), value, true)
}

// Add a optional field to the map.
// The key is a String() Type.
func FieldSO(key string, value Type) MapField {
	return Field(String(key), value, false)
}

// Add a comments
func (f MapField) Comment(c ...string) MapField {
	f.comments = append(f.comments, c...)
	return f
}

func (m *mapType) Match(v any) error {
	mv, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf(notMapFormat, v)
	}

	keys := make(map[string]bool, len(mv))
	for k := range mv {
		keys[k] = true
	}

	errs := make(ErrorSlice, 0, len(mv)+len(keys))
	for _, field := range m.fields {
		key, ok := m.searchKey(&field, keys)
		if !ok && field.required {
			errs = append(errs, fmt.Errorf("not found field for %s", field.fieldKey))
		} else if !ok {
			continue
		} else if err := field.fieldValue.Match(mv[key]); err != nil {
			errs.Append(key, err)
		}
	}

	if len(keys) != 0 && !m.extraFields {
		for key := range keys {
			errs.Append(key, errors.New("extra field"))
		}
	}

	return errs.Return()
}
func (m *mapType) searchKey(field *MapField, keys map[string]bool) (key string, ok bool) {
	for k := range keys {
		if field.fieldKey.Match(k) == nil {
			delete(keys, k)
			return k, true
		}
	}
	return "", false
}

func (m *mapType) HTML(indent string) render.Node {
	if len(m.fields) == 0 {
		if m.extraFields {
			return render.N("", "{...}")
		} else {
			return render.N("", "{}")
		}
	} else if len(m.fields) == 1 && len(m.fields[0].comments) == 0 {
		field := m.fields[0]
		sep := ": "
		if !field.required {
			sep = "?: "
		}
		return render.N("", "{ ",
			field.fieldKey.HTML(indent), sep,
			field.fieldValue.HTML(indent),
			render.If(m.extraFields, func() render.Node { return render.N("", ", ...") }),
			" }")
	}

	indentAdd := indent + "\t"
	return render.N("",
		render.N("", `{`, "\n"),
		render.S(m.fields, "", func(field MapField) render.Node {
			sep := ": "
			if !field.required {
				sep = "?: "
			}
			return render.N("",
				render.S(field.comments, "", func(c string) render.Node {
					return render.N("", indentAdd, render.N("span.sch-comment", "// ", c), "\n")
				}),
				indentAdd,
				field.fieldKey.HTML(indentAdd), sep,
				field.fieldValue.HTML(indentAdd),
				",\n")
		}),
		render.IfS(m.extraFields, render.N("", indentAdd, "...\n")),
		render.N("", indent, `}`),
	)
}
