package sch

import (
	"encoding/json"
	"sniffle/tool/render"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func genHTML(t Type) string {
	s := string(render.Merge(t.HTML("\t")))
	s = strings.TrimPrefix(s, `<!DOCTYPE html>`)
	return s
}

func TestAny(t *testing.T) {
	assert.NoError(t, Any().Match(""))
	assert.NoError(t, Any().Match(1))
	assert.NoError(t, Any().Match(true))
	assert.NoError(t, Any().Match(nil))
	assert.NoError(t, Any().Match(struct{}{}))
	assert.NoError(t, Any().Match(map[string]int{}))
	assert.Equal(t, `<span class=sch-base>...</span>`, genHTML(Any()))
}

func TestAnyBool(t *testing.T) {
	assert.Error(t, AnyBool().Match(""))
	assert.NoError(t, AnyBool().Match(true))
	assert.NoError(t, AnyBool().Match(false))
	assert.Equal(t, `<span class=sch-base>boolean</span>`, genHTML(AnyBool()))
}

func TestBool(t *testing.T) {
	assert.Error(t, True().Match(""))
	assert.Error(t, False().Match(""))
	assert.NoError(t, True().Match(true))
	assert.Error(t, True().Match(false))
	assert.Error(t, False().Match(true))
	assert.NoError(t, False().Match(false))
	assert.Equal(t, `<span class=sch-base>true</span>`, genHTML(True()))
	assert.Equal(t, `<span class=sch-base>false</span>`, genHTML(False()))
}

func TestAsInt(t *testing.T) {
	assert.NoError(t, AsAnyInt().Match(json.Number("-4")))
	assert.NoError(t, AsAnyInt().Match(json.Number("0")))
	assert.NoError(t, AsAnyInt().Match(json.Number("4")))
	assert.Error(t, AsAnyInt().Match(json.Number("-4.0")))
	assert.Error(t, AsAnyInt().Match(""))
	assert.Error(t, AsAnyInt().Match(true))
	assert.Error(t, AsAnyInt().Match(nil))
	assert.Error(t, AsAnyInt().Match(1.5))

	assert.Error(t, AsPositiveInt().Match(json.Number("-1")))
	assert.NoError(t, AsPositiveInt().Match(json.Number("0")))
	assert.NoError(t, AsPositiveInt().Match(json.Number("1")))

	assert.Error(t, AsStrictPositiveInt().Match(json.Number("-1")))
	assert.Error(t, AsStrictPositiveInt().Match(json.Number("0")))
	assert.NoError(t, AsStrictPositiveInt().Match(json.Number("1")))

	assert.NoError(t, AsNegativeInt().Match(json.Number("-1")))
	assert.NoError(t, AsNegativeInt().Match(json.Number("0")))
	assert.Error(t, AsNegativeInt().Match(json.Number("1")))

	assert.NoError(t, AsStrictNegativeInt().Match(json.Number("-1")))
	assert.Error(t, AsStrictNegativeInt().Match(json.Number("0")))
	assert.Error(t, AsStrictNegativeInt().Match(json.Number("1")))

	assert.Error(t, IntervalAsInt(-1, 1).Match(json.Number("-2")))
	assert.NoError(t, IntervalAsInt(-1, 1).Match(json.Number("-1")))
	assert.NoError(t, IntervalAsInt(-1, 1).Match(json.Number("0")))
	assert.NoError(t, IntervalAsInt(-1, 1).Match(json.Number("1")))
	assert.Error(t, IntervalAsInt(-1, 1).Match(json.Number("2")))

	assert.Equal(t, `<span class=sch-int title=Integer>integer</span>`, genHTML(AsAnyInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ 0 ... ]</span>`, genHTML(AsPositiveInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ 1 ... ]</span>`, genHTML(AsStrictPositiveInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ ... 0 ]</span>`, genHTML(AsNegativeInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ ... -1 ]</span>`, genHTML(AsStrictNegativeInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ 3 ... 5 ]</span>`, genHTML(IntervalAsInt(3, 5)))
}

func TestAsStringInt(t *testing.T) {
	assert.NoError(t, AsAnyStringInt().Match("-4"))
	assert.NoError(t, AsAnyStringInt().Match("0"))
	assert.NoError(t, AsAnyStringInt().Match("4"))
	assert.Error(t, AsAnyStringInt().Match("-4.0"))
	assert.Error(t, AsAnyStringInt().Match(""))
	assert.Error(t, AsAnyStringInt().Match(true))
	assert.Error(t, AsAnyStringInt().Match(nil))
	assert.Error(t, AsAnyStringInt().Match(1.5))

	assert.Error(t, AsPositiveStringInt().Match("-1"))
	assert.NoError(t, AsPositiveStringInt().Match("0"))
	assert.NoError(t, AsPositiveStringInt().Match("1"))

	assert.Error(t, AsStrictPositiveStringInt().Match("-1"))
	assert.Error(t, AsStrictPositiveStringInt().Match("0"))
	assert.NoError(t, AsStrictPositiveStringInt().Match("1"))

	assert.NoError(t, AsNegativeStringInt().Match("-1"))
	assert.NoError(t, AsNegativeStringInt().Match("0"))
	assert.Error(t, AsNegativeStringInt().Match("1"))

	assert.NoError(t, AsStrictNegativeStringInt().Match("-1"))
	assert.Error(t, AsStrictNegativeStringInt().Match("0"))
	assert.Error(t, AsStrictNegativeStringInt().Match("1"))

	assert.Equal(t, `<span class=sch-int title="Integer into a string">string(integer)</span>`, genHTML(AsAnyStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ 0 ... ])</span>`, genHTML(AsPositiveStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ 1 ... ])</span>`, genHTML(AsStrictPositiveStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ ... 0 ])</span>`, genHTML(AsNegativeStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ ... -1 ])</span>`, genHTML(AsStrictNegativeStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ 3 ... 5 ])</span>`, genHTML(IntervalAsStringInt(3, 5)))

	assert.Equal(t, `string(integer)`, AsAnyStringInt().String())
	assert.Equal(t, `string([ 3 ... 5 ])`, IntervalAsStringInt(3, 5).String())
}

func TestNil(t *testing.T) {
	assert.NoError(t, Nil().Match(nil))
	assert.Error(t, Nil().Match(1))
	assert.Equal(t, `<span class=sch-base>null</span>`, genHTML(Nil()))
}

func TestAnyString(t *testing.T) {
	assert.NoError(t, AnyString().Match(""))
	assert.NoError(t, AnyString().Match("abc"))
	assert.Error(t, AnyString().Match(1))
	assert.Equal(t, `<span class=sch-base>string</span>`, genHTML(AnyString()))
}

func TestNotEmptyString(t *testing.T) {
	assert.NoError(t, NotEmptyString().Match("abc"))
	assert.Error(t, NotEmptyString().Match(""))
	assert.Error(t, NotEmptyString().Match(1))
	assert.Equal(t, `<span class=sch-base>not-empty-string</span>`, genHTML(NotEmptyString()))
}

func TestArroundString(t *testing.T) {
	assert.NoError(t, ArroundString("abc\t").Match("abc"))
	assert.NoError(t, ArroundString("abc").Match(" abc "))
	assert.NoError(t, ArroundString("\nabc").Match("ABC"))
	assert.Error(t, ArroundString("abc").Match(""))
	assert.Error(t, ArroundString("abc").Match(1))
	assert.Equal(t, `<span class=sch-str>~&#34;abc&#34;</span>`, genHTML(ArroundString(" aBc ")))
	assert.Equal(t, "~abc", ArroundString("abc").String())
}

func TestString(t *testing.T) {
	assert.NoError(t, String("abc").Match("abc"))
	assert.Error(t, String("abc").Match(" abc "))
	assert.Error(t, String("abc").Match("ABC"))
	assert.Error(t, String("abc").Match(""))
	assert.Error(t, String("abc").Match(1))
	assert.Equal(t, `<span class=sch-str>&#34;abC&#34;</span>`, genHTML(String("abC")))
	assert.Equal(t, `abC`, String("abC").String())
}
