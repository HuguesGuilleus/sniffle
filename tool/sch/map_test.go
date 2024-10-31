package sch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	m := NewMap().
		FieldSR("k1", False()).
		FieldSO("k2", True())

	assert.NoError(t, m.Match(map[string]any{"k1": false}))
	assert.NoError(t, m.Match(map[string]any{"k1": false, "k2": true}))
	assert.Error(t, m.Match(map[string]any{"k1": 1}))
	assert.Error(t, m.Match(map[string]any{}))
	assert.Error(t, m.Match(map[string]any{"k1": false, "x": nil}))
	assert.Error(t, m.Match(1))

	m.ExtraFields()
	assert.NoError(t, m.Match(map[string]any{"k1": false}))
	assert.NoError(t, m.Match(map[string]any{"k1": false, "x": nil}))

	assert.Equal(t, "{"+
		"\n\t\t<span class=sch-str>&#34;k1&#34;</span>: <span class=sch-base>false</span>,"+
		"\n\t\t<span class=sch-str>&#34;k2&#34;</span>?: <span class=sch-base>true</span>,"+
		"\n\t}", genHTML(m))
}
