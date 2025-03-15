// sch check that data observe a schema.
package sch

import "sniffle/tool/render"

type defType struct {
	id string
	ty Type
}

// Wrap the type into a definition.
// The returned type.Match() use t.Match().
// For type.HTML(), return only `<id>`.
func Def(id string, t Type) (Type, render.Node) {
	return defType{id, t}, render.N("", render.Na("span.sch-def", "id", "sch.def."+id).N("type ", id, " = ", t.HTML("")))
}

func (t defType) Match(v any) error { return t.ty.Match(v) }

func (t defType) HTML(_ string) render.Node {
	return render.Na("a.sch-def", "href", "#sch.def."+t.id).N(render.H("type&nbsp;"), t.id)
}
