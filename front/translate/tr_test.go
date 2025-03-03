package translate

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"sniffle/tool/language"
	"sniffle/tool/render"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	tr := make([]map[string]string, len(Langs))
	m := make(map[string]struct{})
	for i, l := range Langs {
		tr[i] = make(map[string]string)
		collect(tr[i], "T", reflect.ValueOf(T[l]))
		mergeMap(m, tr[i])
	}

	os.WriteFile("tr.html", render.Merge(render.N("html",
		render.N("head",
			render.N("title", "Translate"),
			render.N("style", ``+
				`body{font-family:sans;font-size:large}`+
				`table{width:100%;border-collapse:collapse}`+
				`thead{position:sticky;top:0}`+
				`td,th{border:solid .3ex grey;background:#eee}`+
				`th{background:#ddd}`+
				`.empty{background:grey}`+
				`tr:hover *{background:#bbb}`),
		),
		render.N("body",
			render.S(Langs, "; ", func(l language.Language) render.Node {
				return render.N("",
					render.Na("label", "for", "lang-"+l.String()).N(l.Human()),
					render.Na("input", "type", "checkbox").
						A("id", "lang-"+l.String()).
						A("data-i", l.String()).
						N(),
				)
			}),
			render.N("br"), render.N("br"),
			render.Map(namespaceMap(m), func(ns string, keys []string) render.Node {
				return render.N("",
					render.N("h1", ns),
					render.N("table",
						render.N("thead", render.N("tr",
							render.N("th", "K"),
							render.S(Langs, "", func(l language.Language) render.Node {
								return render.Na("th", "data-l", l.String()).N(l.Human())
							}),
						)),
						render.S(keys, "", func(key string) render.Node {
							return render.N("tr",
								render.N("td", key),
								render.S2(tr, "", func(i int, m map[string]string) render.Node {
									if m[key] == "" {
										return render.Na("td.empty", "data-l", Langs[i].String()).N()
									}
									return render.Na("td", "data-l", Langs[i].String()).N(m[key])
								}),
							)
						}),
					),
				)
			}),
			render.N("script", render.H(``+
				`const qsa=(q,f)=>document.querySelectorAll(q).forEach(f),`+
				`f=i=>qsa("[data-l="+i.dataset.i+"]",e=>e.hidden=!i.checked);`+
				`qsa("input",i=>[f(i),i.onchange=_=>f(i)])`,
			)),
		),
	)), 0o664)
}

func collect(m map[string]string, base string, v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		m[base] = v.String()
	case reflect.Array, reflect.Slice:
		for i := range v.Len() {
			collect(m, fmt.Sprintf("%s[%02d]", base, i), v.Index(i))
		}
	case reflect.Map:
		for iter := v.MapRange(); iter.Next(); {
			collect(m, base+"["+iter.Key().String()+"]", iter.Value())
		}
	case reflect.Struct:
		for i := range v.NumField() {
			field := v.Type().Field(i)
			collect(m, base+"."+field.Name, v.Field(i))
		}
	default:
		panic("unknown kind type: " + v.Kind().String())
	}
}

func mergeMap(out map[string]struct{}, src map[string]string) {
	for k := range src {
		out[k] = struct{}{}
	}
}

func namespaceMap(m map[string]struct{}) map[string][]string {
	out := make(map[string][]string)
	for key := range m {
		ns := strings.FieldsFunc(key, func(r rune) bool { return r == '.' || r == '[' })[1]
		out[ns] = append(out[ns], key)
	}
	for _, keys := range out {
		slices.Sort(keys)
	}
	return out
}
