package sch

import (
	"errors"
	"fmt"
	"sniffle/tool/render"
)

type Map struct {
	fields      []mapField
	extraFields bool
}
type mapField struct {
	fieldKey   TypeStringer
	fieldValue Type
	required   bool
}

func NewMap() *Map { return new(Map) }

// Add a field to the map.
func (m *Map) Field(key TypeStringer, value Type, required bool) *Map {
	m.fields = append(m.fields, mapField{key, value, required})
	return m
}

// Add a required field to the map.
// The key is a String() Type.
func (m *Map) FieldSR(key string, value Type) *Map {
	return m.Field(String(key), value, true)
}

// Add a optional field to the map.
// The key is a String() Type.
func (m *Map) FieldSO(key string, value Type) *Map {
	return m.Field(String(key), value, false)
}

func (m *Map) ExtraFields() *Map {
	m.extraFields = true
	return m
}

func (m *Map) Match(v any) error {
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
			errs = append(errs, fmt.Errorf("not found key for %s", field.fieldKey))
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
func (m *Map) searchKey(field *mapField, keys map[string]bool) (key string, ok bool) {
	for k := range keys {
		if field.fieldKey.Match(k) == nil {
			delete(keys, k)
			return k, true
		}
	}
	return "", false
}

func (m *Map) HTML(indent string) render.Node {
	indentAdd := indent + "\t"
	return render.N("",
		render.N("", `{`, "\n"),
		render.S(m.fields, "", func(item mapField) render.Node {
			sep := ": "
			if !item.required {
				sep = "?: "
			}
			return render.N("", indentAdd,
				item.fieldKey.HTML(indentAdd), sep,
				item.fieldValue.HTML(indentAdd),
				",\n")
		}),
		render.N("", indent, `}`),
	)
}
