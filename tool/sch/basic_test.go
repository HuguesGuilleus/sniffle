package sch

import (
	"encoding/json"
	"io"
	"os"
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

func TestFloat(t *testing.T) {
	assert.NoError(t, AnyFloat().Match(json.Number("-4.2")))
	assert.NoError(t, AnyFloat().Match(json.Number("00.0")))
	assert.NoError(t, AnyFloat().Match(json.Number("00")))
	assert.NoError(t, AnyFloat().Match(json.Number("4.3")))
	assert.NoError(t, AnyFloat().Match(json.Number("4.0")))
	assert.NoError(t, AnyFloat().Match(json.Number("4")))
	assert.NoError(t, AnyFloat().Match(json.Number("04")))
	assert.Error(t, AnyFloat().Match(json.Number("j")))
	assert.Error(t, AnyFloat().Match(""))
	assert.Error(t, AnyFloat().Match(true))
	assert.Error(t, AnyFloat().Match(nil))

	assert.Error(t, PositiveFloat().Match(json.Number("-1.0")))
	assert.NoError(t, PositiveFloat().Match(json.Number("0.0")))
	assert.NoError(t, PositiveFloat().Match(json.Number("0.1")))

	assert.Error(t, StrictPositiveFloat().Match(json.Number("-1.0")))
	assert.Error(t, StrictPositiveFloat().Match(json.Number("0.0")))
	assert.NoError(t, StrictPositiveFloat().Match(json.Number("0.1")))

	assert.NoError(t, NegativeFloat().Match(json.Number("-1.0")))
	assert.NoError(t, NegativeFloat().Match(json.Number("0.0")))
	assert.Error(t, NegativeFloat().Match(json.Number("0.1")))

	assert.NoError(t, StrictNegativeFloat().Match(json.Number("-0.1")))
	assert.Error(t, StrictNegativeFloat().Match(json.Number("0.0")))
	assert.Error(t, StrictNegativeFloat().Match(json.Number("1.0")))

	assert.Error(t, IntervalFloat(-1, 1).Match(json.Number("-2.0")))
	assert.NoError(t, IntervalFloat(-1, 1).Match(json.Number("-1.0")))
	assert.NoError(t, IntervalFloat(-1, 1).Match(json.Number("0.0")))
	assert.NoError(t, IntervalFloat(-1, 1).Match(json.Number("1.0")))
	assert.Error(t, IntervalFloat(-1, 1).Match(json.Number("2.0")))

	assert.Equal(t, `<span class=sch-float title=Float64>float64</span>`, genHTML(AnyFloat()))
	assert.Equal(t, `<span class=sch-float title=Float64>( 0.0 .. )</span>`, genHTML(PositiveFloat()))
	assert.Equal(t, `<span class=sch-float title=Float64>( 0.0 > .. )</span>`, genHTML(StrictPositiveFloat()))
	assert.Equal(t, `<span class=sch-float title=Float64>( .. 0.0 )</span>`, genHTML(NegativeFloat()))
	assert.Equal(t, `<span class=sch-float title=Float64>( .. < 0.0 )</span>`, genHTML(StrictNegativeFloat()))
	assert.Equal(t, `<span class=sch-float title=Float64>( -1.567 .. 2.3 )</span>`, genHTML(IntervalFloat(-1.567, 2.3)))
	assert.Equal(t, `<span class=sch-float title=Float64>( -1.1 .. 2.3 )</span>`, genHTML(IntervalFloat(-1.1, 2.3)))
	assert.Equal(t, `<span class=sch-float title=Float64>3.123456789</span>`, genHTML(ConstFloat(3.123456789)))
	assert.Equal(t, `<span class=sch-float title=Float64>3.0</span>`, genHTML(ConstFloat(3)))
	assert.Equal(t, `<span class=sch-float title=Float64>0.0</span>`, genHTML(ConstFloat(0)))
}

func TestStringFloat(t *testing.T) {
	assert.NoError(t, AnyStringFloat().Match("-4.2"))
	assert.NoError(t, AnyStringFloat().Match("00.0"))
	assert.NoError(t, AnyStringFloat().Match("00"))
	assert.NoError(t, AnyStringFloat().Match("4.3"))
	assert.NoError(t, AnyStringFloat().Match("4.0"))
	assert.NoError(t, AnyStringFloat().Match("4"))
	assert.NoError(t, AnyStringFloat().Match("04"))
	assert.Error(t, AnyStringFloat().Match("j"))
	assert.Error(t, AnyStringFloat().Match(""))
	assert.Error(t, AnyStringFloat().Match(true))
	assert.Error(t, AnyStringFloat().Match(nil))

	assert.Error(t, PositiveStringFloat().Match("-1.0"))
	assert.NoError(t, PositiveStringFloat().Match("0.0"))
	assert.NoError(t, PositiveStringFloat().Match("0.1"))

	assert.Error(t, StrictPositiveStringFloat().Match("-1.0"))
	assert.Error(t, StrictPositiveStringFloat().Match("0.0"))
	assert.NoError(t, StrictPositiveStringFloat().Match("0.1"))

	assert.NoError(t, NegativeStringFloat().Match("-1.0"))
	assert.NoError(t, NegativeStringFloat().Match("0.0"))
	assert.Error(t, NegativeStringFloat().Match("0.1"))

	assert.NoError(t, StrictNegativeStringFloat().Match("-0.1"))
	assert.Error(t, StrictNegativeStringFloat().Match("0.0"))
	assert.Error(t, StrictNegativeStringFloat().Match("1.0"))

	assert.Error(t, IntervalStringFloat(-1, 1).Match("-2.0"))
	assert.NoError(t, IntervalStringFloat(-1, 1).Match("-1.0"))
	assert.NoError(t, IntervalStringFloat(-1, 1).Match("0.0"))
	assert.NoError(t, IntervalStringFloat(-1, 1).Match("1.0"))
	assert.Error(t, IntervalStringFloat(-1, 1).Match("2.0"))

	assert.Equal(t, `<span class=sch-float title="Float64 into a string">string(3.0)</span>`, genHTML(ConstStringFloat(3)))
	assert.Equal(t, `string(3.0)`, ConstStringFloat(3).String())
}

func TestAsInt(t *testing.T) {
	assert.NoError(t, AnyInt().Match(json.Number("-4")))
	assert.NoError(t, AnyInt().Match(json.Number("0")))
	assert.NoError(t, AnyInt().Match(json.Number("4")))
	assert.Error(t, AnyInt().Match(json.Number("-4.0")))
	assert.Error(t, AnyInt().Match(""))
	assert.Error(t, AnyInt().Match(true))
	assert.Error(t, AnyInt().Match(nil))
	assert.Error(t, AnyInt().Match(1.5))

	assert.Error(t, PositiveInt().Match(json.Number("-1")))
	assert.NoError(t, PositiveInt().Match(json.Number("0")))
	assert.NoError(t, PositiveInt().Match(json.Number("1")))

	assert.Error(t, StrictPositiveInt().Match(json.Number("-1")))
	assert.Error(t, StrictPositiveInt().Match(json.Number("0")))
	assert.NoError(t, StrictPositiveInt().Match(json.Number("1")))

	assert.NoError(t, NegativeInt().Match(json.Number("-1")))
	assert.NoError(t, NegativeInt().Match(json.Number("0")))
	assert.Error(t, NegativeInt().Match(json.Number("1")))

	assert.NoError(t, StrictNegativeInt().Match(json.Number("-1")))
	assert.Error(t, StrictNegativeInt().Match(json.Number("0")))
	assert.Error(t, StrictNegativeInt().Match(json.Number("1")))

	assert.Error(t, IntervalInt(-1, 1).Match(json.Number("-2")))
	assert.NoError(t, IntervalInt(-1, 1).Match(json.Number("-1")))
	assert.NoError(t, IntervalInt(-1, 1).Match(json.Number("0")))
	assert.NoError(t, IntervalInt(-1, 1).Match(json.Number("1")))
	assert.Error(t, IntervalInt(-1, 1).Match(json.Number("2")))

	assert.Error(t, ConstInt(1).Match(json.Number("-2")))
	assert.Error(t, ConstInt(1).Match(json.Number("-1")))
	assert.Error(t, ConstInt(1).Match(json.Number("0")))
	assert.NoError(t, ConstInt(1).Match(json.Number("1")))
	assert.Error(t, ConstInt(1).Match(json.Number("2")))

	assert.Equal(t, `<span class=sch-int title=Integer>integer</span>`, genHTML(AnyInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ 0 .. ]</span>`, genHTML(PositiveInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ 1 .. ]</span>`, genHTML(StrictPositiveInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ .. 0 ]</span>`, genHTML(NegativeInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ .. -1 ]</span>`, genHTML(StrictNegativeInt()))
	assert.Equal(t, `<span class=sch-int title=Integer>[ 3 .. 5 ]</span>`, genHTML(IntervalInt(3, 5)))
	assert.Equal(t, `<span class=sch-int title=Integer>3</span>`, genHTML(ConstInt(3)))
}

func TestAsStringInt(t *testing.T) {
	assert.NoError(t, AnyStringInt().Match("-4"))
	assert.NoError(t, AnyStringInt().Match("0"))
	assert.NoError(t, AnyStringInt().Match("4"))
	assert.Error(t, AnyStringInt().Match("-4.0"))
	assert.Error(t, AnyStringInt().Match(""))
	assert.Error(t, AnyStringInt().Match(true))
	assert.Error(t, AnyStringInt().Match(nil))
	assert.Error(t, AnyStringInt().Match(1.5))

	assert.Error(t, PositiveStringInt().Match("-1"))
	assert.NoError(t, PositiveStringInt().Match("0"))
	assert.NoError(t, PositiveStringInt().Match("1"))

	assert.Error(t, StrictPositiveStringInt().Match("-1"))
	assert.Error(t, StrictPositiveStringInt().Match("0"))
	assert.NoError(t, StrictPositiveStringInt().Match("1"))

	assert.NoError(t, NegativeStringInt().Match("-1"))
	assert.NoError(t, NegativeStringInt().Match("0"))
	assert.Error(t, NegativeStringInt().Match("1"))

	assert.NoError(t, StrictNegativeStringInt().Match("-1"))
	assert.Error(t, StrictNegativeStringInt().Match("0"))
	assert.Error(t, StrictNegativeStringInt().Match("1"))

	assert.Equal(t, `<span class=sch-int title="Integer into a string">string(integer)</span>`, genHTML(AnyStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ 0 .. ])</span>`, genHTML(PositiveStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ 1 .. ])</span>`, genHTML(StrictPositiveStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ .. 0 ])</span>`, genHTML(NegativeStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ .. -1 ])</span>`, genHTML(StrictNegativeStringInt()))
	assert.Equal(t, `<span class=sch-int title="Integer into a string">string([ 3 .. 5 ])</span>`, genHTML(IntervalStringInt(3, 5)))

	assert.Equal(t, `string(integer)`, AnyStringInt().String())
	assert.Equal(t, `string([ 3 .. 5 ])`, IntervalStringInt(3, 5).String())
}

func TestPrint(t *testing.T) {
	p := Print("test")
	assert.Equal(t, `<span class=sch-base>[log:test]</span>`, genHTML(p))

	defer func(stdout *os.File) { os.Stdout = stdout }(os.Stdout)
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w

	assert.NoError(t, p.Match(1))
	w.Close()

	out, err := io.ReadAll(r)
	assert.NoError(t, err)
	assert.EqualValues(t, "[test] 1\n", out)
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
