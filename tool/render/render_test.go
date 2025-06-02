package render

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	h := Merge(Na("body.cl1.cl2.cl3", "hidden", "").
		A("data-toescape", "<must' &escape\">").
		A("data-normal", "yolo/").
		N(
			"Hello ",
			N("i.cl", "World"),
			Na("img", "src", "icon.png"),
			H("&#33;"),
			"\n",
			true,
			uint(42),
			-1,
			[]H{"a", "b"},
			time.Second,
			nil,
		))
	assert.Equal(t, `<!DOCTYPE html>`+
		`<body class="cl1 cl2 cl3"`+
		` hidden`+
		` data-toescape="<must' &amp;escape&#34;>"`+
		` data-normal=yolo/`+
		` >`+
		`Hello `+
		`<i class=cl>World</i>`+
		`<img src=icon.png>`+
		`&#33;`+
		"\n"+
		`true`+
		`42`+
		`-1`+
		`ab`+
		`1s`+
		`</body>`,
		string(h))
}

func TestZ(t *testing.T) {
	assert.Nil(t, Z.mergeSlice(nil))
}
func TestNoMarkup(t *testing.T) {
	n := N("", 1)
	assert.Equal(t, `1`, string(n.mergeSlice(nil)))
}

func TestIf(t *testing.T) {
	n := If(false, func() Node { return N("img") })
	assert.Nil(t, n.mergeSlice(nil))
	n = If(true, func() Node { return N("img") })
	assert.Equal(t, `<img>`, string(n.mergeSlice(nil)))
}
func TestIfS(t *testing.T) {
	n := IfS(false, N("img"))
	assert.Nil(t, n.mergeSlice(nil))
	n = IfS(true, N("img"))
	assert.Equal(t, `<img>`, string(n.mergeSlice(nil)))
}
func TestIfElse(t *testing.T) {
	n := IfElse(true, func() Node { return N("a") }, func() Node { return N("b") })
	assert.Equal(t, `<a></a>`, string(n.mergeSlice(nil)))
	n = IfElse(false, func() Node { return N("a") }, func() Node { return N("b") })
	assert.Equal(t, `<b></b>`, string(n.mergeSlice(nil)))
}

func TestIsSAny(t *testing.T) {
	assert.Equal(t, 3, IfSAny(false, 1, 3))
	assert.Equal(t, 1, IfSAny(true, 1, 3))
}

func TestEspaceString(t *testing.T) {
	n := N("", `"Fran & Freddie's Diner" <tasty@example.com> `)
	assert.Equal(t, `&#34;Fran &amp; Freddie&#39;s Diner&#34; &lt;tasty@example.com&gt; `, string(n.mergeSlice(nil)))
}
func TestTime(t *testing.T) {
	assert.Equal(t, `<time datetime=2024-02-14T20:21:22Z>2024-02-14 20:21:22 UTC</time>`,
		string(renderChild(nil, time.Date(2024, time.February, 14, 20, 21, 32, 1, time.FixedZone("TEST", 10)))))
}
func TestDate(t *testing.T) {
	assert.Equal(t, `<time datetime=2024-02-14>2024-02-14</time>`,
		string(renderChild(nil, time.Date(2024, time.February, 14, 0, 0, 0, 0, DateZone))))
}
func TestShortDate(t *testing.T) {
	assert.Equal(t, `<time datetime=2024-02-14>2024_02_14</time>`,
		string(renderChild(nil, time.Date(2024, time.February, 14, 0, 0, 0, 0, ShortDateZone))))
}

func TestIntType(t *testing.T) {
	assert.Equal(t, `-123456789`, string(renderChild(nil, Int(-123_456_789))))
}
func TestNumber(t *testing.T) {
	assert.Equal(t, "-123\u202F456\u202F789", string(renderChild(nil, -123_456_789)))
	assert.Equal(t, "-23\u202F456\u202F789", string(renderChild(nil, -23_456_789)))
	assert.Equal(t, "-3\u202F456\u202F789", string(renderChild(nil, -3_456_789)))
	assert.Equal(t, "-3\u202F056\u202F789", string(renderChild(nil, -3_056_789)))
	assert.Equal(t, "1\u202F000", string(renderChild(nil, 1_000)))
	assert.Equal(t, `-789`, string(renderChild(nil, -789)))
	assert.Equal(t, `-89`, string(renderChild(nil, -89)))
	assert.Equal(t, `-9`, string(renderChild(nil, -9)))
	assert.Equal(t, `0`, string(renderChild(nil, 0)))
	assert.Equal(t, `100`, string(renderChild(nil, 100)))
	assert.Equal(t, `10`, string(renderChild(nil, 10)))

	for i := range 1_111_111 {
		assert.Equal(t,
			strings.ReplaceAll(string(renderUint64(nil, uint64(i))), "\u202F", ""),
			strconv.Itoa(i), i)
	}
}

func TestArray(t *testing.T) {
	n := N("_", []any{1, "A", true})
	assert.Equal(t, `<_>1Atrue</_>`, string(n.mergeSlice(nil)))
}

func TestEmptyNode(t *testing.T) {
	n := N("link")
	assert.Equal(t, `<link>`, string(n.mergeSlice(nil)))
}

func TestMap(t *testing.T) {
	m := map[int]bool{2: false, 1: true}
	h := Merge(N("ul", Map(m, func(k int, v bool) Node {
		return N("li", "k:", k, " => v:", v)
	})))
	assert.Equal(t, `<!DOCTYPE html>`+
		`<ul>`+
		`<li>k:1 =&gt; v:true</li>`+
		`<li>k:2 =&gt; v:false</li>`+
		`</ul>`,
		string(h))
}
func TestMapReverse(t *testing.T) {
	m := map[int]bool{2: false, 1: true}
	h := Merge(N("ul", MapReverse(m, func(k int, v bool) Node {
		return N("li", "k:", k, " => v:", v)
	})))
	assert.Equal(t, `<!DOCTYPE html>`+
		`<ul>`+
		`<li>k:2 =&gt; v:false</li>`+
		`<li>k:1 =&gt; v:true</li>`+
		`</ul>`,
		string(h))
}

func TestSlice(t *testing.T) {
	s := []bool{true, false}
	assert.Equal(t,
		[]Node{
			{"code", nil, []any{true}},
			{"", nil, []any{H(",")}},
			{"code", nil, []any{false}},
		},
		S(s, ",", func(b bool) Node {
			return N("code", b)
		}))
	assert.Nil(t, S[any](nil, "", nil))
}

func TestSliceSeparator(t *testing.T) {
	s := []bool{true, false}
	assert.Equal(t,
		[]Node{
			{"code", nil, []any{0, true}},
			{"", nil, []any{H("/")}},
			{"code", nil, []any{1, false}},
		},
		S2(s, "/", func(i int, b bool) Node {
			return N("code", i, b)
		}))
	assert.Nil(t, S2[any](nil, "", nil))
}
