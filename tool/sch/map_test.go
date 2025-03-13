package sch

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assertCall := 0
	assertThis := map[string]any{"k1": false, "k2": true, "as": "UP"}
	m := Map(
		FieldSR("k1", False()),
		FieldSO("k2", True()),
		FieldSO("as", AnyString()).Assert(`uppercase`, func(this map[string]any, field any) error {
			assertCall++
			if f := field.(string); strings.ToUpper(f) != f {
				return fmt.Errorf("need upper case string")
			} else if reflect.ValueOf(assertThis).UnsafePointer() != reflect.ValueOf(this).UnsafePointer() {
				return fmt.Errorf("not asserThis")
			}
			return nil
		}),
	)

	assert.NoError(t, m.Match(map[string]any{"k1": false}))
	assert.NoError(t, m.Match(map[string]any{"k1": false, "k2": true}))
	assert.NoError(t, m.Match(assertThis))
	assert.Equal(t, 1, assertCall)
	assert.Error(t, m.Match(map[string]any{"k1": 1}))
	assert.Error(t, m.Match(map[string]any{}))
	assert.Error(t, m.Match(map[string]any{"k1": false, "x": nil}))
	assert.Error(t, m.Match(1))
	assert.Error(t, m.Match(map[string]any{"k1": false, "k2": true, "as": "UP"}))
	assert.Error(t, m.Match(map[string]any{"k1": false, "k2": true, "as": 1}))
	assertThis["as"] = "up"
	assert.Error(t, m.Match(assertThis))
	assert.Equal(t, 3, assertCall)

	assert.Error(t, Map(Assert(``, func(this map[string]any, _ any) error { return errors.New("e") })).Match(map[string]any{}))

	m = MapExtra(
		FieldSR("k1", False()),
		FieldSO("k2", True()).
			Comment("c1", "c2").
			Assert(`a1`, func(this map[string]any, field any) error { return nil }).
			Assert(`a2`, func(this map[string]any, field any) error { return nil }),
		Assert(`alone`, func(this map[string]any, _ any) error { return nil }),
	)
	assert.NoError(t, m.Match(map[string]any{"k1": false}))
	assert.NoError(t, m.Match(map[string]any{"k1": false, "x": nil}))

	assert.Equal(t, "{"+
		"\n\t\t<span class=sch-str>&#34;k1&#34;</span>: <span class=sch-base>false</span>,"+
		"\n\t\t<span class=sch-comment>// c1</span>"+
		"\n\t\t<span class=sch-comment>// c2</span>"+
		"\n\t\t<span class=sch-str>&#34;k2&#34;</span>?: <span class=sch-base>true</span>"+
		" #assert <span class=sch-assert>a1</span>"+
		" #assert <span class=sch-assert>a2</span>"+
		","+
		"\n\t\t#assert <span class=sch-assert>alone</span>"+
		"\n\t\t..."+
		"\n\t}", genHTML(m))

	assert.Equal(t, "{\n\t\t<span class=sch-str>&#34;k1&#34;</span>: <span class=sch-base>false</span>,\n\t}",
		genHTML(Map(FieldSR("k1", False()))),
	)
	assert.Equal(t, "{\n\t\t<span class=sch-str>&#34;k1&#34;</span>?: <span class=sch-base>false</span>,\n\t\t...\n\t}",
		genHTML(MapExtra(FieldSO("k1", False()))),
	)

	assert.Equal(t, "{}", genHTML(Map()))
	assert.Equal(t, "{...}", genHTML(MapExtra()))

	assert.Equal(t, "{\n\n\t}", genHTML(Map(BlankField())))
}
