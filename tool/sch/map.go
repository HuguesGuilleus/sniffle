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
	asserts    []assertFunc
}
type assertFunc struct {
	name string
	test func(this map[string]any, field any) error
}

func Map(fields ...MapField) Type      { return &mapType{fields, false} }
func MapExtra(fields ...MapField) Type { return &mapType{fields, true} }

// Add a black line
func BlankField() MapField {
	return Field(nil, nil, false)
}

// Add a field to the map.
func Field(key TypeStringer, value Type, required bool) MapField {
	return MapField{key, value, required, nil, nil}
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

func Assert(name string, test func(this map[string]any, _ any) error) MapField {
	return MapField{nil, nil, false, nil, []assertFunc{{name, test}}}
}

// TODO: add in method
func (f MapField) Assert(name string, test func(this map[string]any, field any) error) MapField {
	f.asserts = append(f.asserts, assertFunc{name, test})
	return f
}

// Comment appends some comment line to this field.
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
		if field.fieldKey == nil {
			if len(errs) == 0 {
				for _, assert := range field.asserts {
					if err := assert.test(mv, nil); err != nil {
						errs = append(errs, err)
					}
				}
			}
			continue
		}

		key, ok := m.searchKey(&field, keys)
		if !ok && field.required {
			errs = append(errs, fmt.Errorf("not found field for %s", field.fieldKey))
			continue
		} else if !ok {
			continue
		}
		fieldValue := mv[key]
		if err := field.fieldValue.Match(fieldValue); err != nil {
			errs.Append(key, err)
		} else {
			for _, assert := range field.asserts {
				errs.Append(key, assert.test(mv, fieldValue))
			}
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
	}

	indentAdd := indent + "\t"
	return render.N("",
		render.N("", `{`, "\n"),
		render.S(m.fields, "", func(field MapField) render.Node {
			comments := render.S(field.comments, "", func(c string) render.Node {
				return render.N("", indentAdd, render.N("span.sch-comment", "// ", c), "\n")
			})
			asserts := render.S(field.asserts, " ", func(a assertFunc) render.Node {
				return render.N("", "#assert ", render.N("span.sch-assert", a.name))
			})
			if field.fieldKey == nil {
				if len(asserts) == 0 {
					return render.N("", comments, "\n")
				}
				return render.N("",
					comments,
					indentAdd, asserts, "\n",
				)
			} else {
				sep := ": "
				if !field.required {
					sep = "?: "
				}
				return render.N("",
					comments,
					render.If(field.fieldKey != nil, func() render.Node {
						return render.N("", indentAdd,
							field.fieldKey.HTML(indentAdd), sep,
							field.fieldValue.HTML(indentAdd),
						)
					}),
					render.IfS(len(asserts) > 0, render.N("", " ", asserts)),
					",\n",
				)
			}
		}),
		render.IfS(m.extraFields, render.N("", indentAdd, "...\n")),
		render.N("", indent, `}`),
	)
}
