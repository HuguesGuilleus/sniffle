package sch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnd(t *testing.T) {
	and := And(String("b"), ArroundString("B"))
	assert.NoError(t, and.Match("b"))
	assert.Error(t, and.Match("B"))
	assert.Error(t, and.Match(1))
	assert.Equal(t, `<span class=sch-str>&#34;b&#34;</span> & <span class=sch-str>~&#34;b&#34;</span>`, genHTML(and))
}

func TestOr(t *testing.T) {
	or := Or(String("a"), True())
	assert.NoError(t, or.Match("a"))
	assert.NoError(t, or.Match(true))
	assert.Error(t, or.Match(false))
	assert.Equal(t, `<span class=sch-str>&#34;a&#34;</span> | <span class=sch-base>true</span>`, genHTML(or))
}

func TestEnumString(t *testing.T) {
	assert.Equal(t, enumString{orType{String("a"), String("b")}, "a|b"}, EnumString("a", "b"))
	assert.Equal(t, "a|b", EnumString("a", "b").String())
}
