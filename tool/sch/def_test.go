package sch

import (
	"sniffle/tool/render"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDef(t *testing.T) {
	d, defHTML := Def("d", Map())

	assert.NoError(t, d.Match(map[string]any{}))
	assert.Error(t, d.Match(1))

	assert.Equal(t, `<a class=sch-def href=#sch.def.d>type&nbsp;d</a>`, genHTML(d))
	assert.Equal(t, `<!DOCTYPE html>`+
		`<span class=sch-def id=sch.def.d>type d = {}</span>`,
		string(render.Merge(defHTML)))
}
