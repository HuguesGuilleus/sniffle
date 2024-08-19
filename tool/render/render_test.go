package render

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	h := Merge(No("body.cl1.cl2.cl3",
		A("hidden", "").
			A("data-toescape", "<must' &escape\">").
			A("data-normal", "yolo/"),
		"Hello ",
		N("i.cl", "World"),
		H("&#33;"),
		"\n",
		true,
		uint(42),
		-1,
		boolHTML(true),
		[]H{"a", "b"},
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
		`&#33;`+
		"\n"+
		`true`+
		`42`+
		`-1`+
		`TRUE`+
		`ab`+
		`</body>`,
		string(h))
}

func TestZ(t *testing.T) {
	assert.Nil(t, Z.mergeSlice(nil))
}
func TestExclamation(t *testing.T) {
	n := N("", 1)
	assert.Equal(t, `1`, string(n.mergeSlice(nil)))
}

func TestIf(t *testing.T) {
	n := If(false, func() Node { return N("img") })
	assert.Nil(t, n.mergeSlice(nil))
	n = If(true, func() Node { return N("img") })
	assert.Equal(t, `<img>`, string(n.mergeSlice(nil)))
}
func TestIfElse(t *testing.T) {
	n := IfElse(true, func() Node { return N("a") }, func() Node { return N("b") })
	assert.Equal(t, `<a></a>`, string(n.mergeSlice(nil)))
	n = IfElse(false, func() Node { return N("a") }, func() Node { return N("b") })
	assert.Equal(t, `<b></b>`, string(n.mergeSlice(nil)))
}

func TestTime(t *testing.T) {
	assert.Equal(t, `<time datetime=2024-02-14T20:21:22Z>2024-02-14 20:21:22 UTC</time>`,
		string(renderChild(nil, time.Date(2024, time.February, 14, 20, 21, 32, 1, time.FixedZone("TEST", 10)))))
}
func TestDate(t *testing.T) {
	assert.Equal(t, `<time datetime=2024-02-14>2024-02-14</time>`,
		string(renderChild(nil, time.Date(2024, time.February, 14, 0, 0, 0, 0, DateZone))))
}

func TestIntType(t *testing.T) {
	assert.Equal(t, `-123456789`, string(renderChild(nil, Int(-123_456_789))))
}
func TestNumber(t *testing.T) {
	assert.Equal(t, `-123 456 789`, string(renderChild(nil, -123_456_789)))
}

func TestArray(t *testing.T) {
	n := N("_", []any{1, "A", true})
	assert.Equal(t, `<_>1Atrue</_>`, string(n.mergeSlice(nil)))
}

func TestEmptyNode(t *testing.T) {
	n := N("link")
	assert.Equal(t, `<link>`, string(n.mergeSlice(nil)))
}

type boolHTML bool

func (b boolHTML) HTML() H {
	if b {
		return "TRUE"
	} else {
		return "FALSE"
	}
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

func TestSlice(t *testing.T) {
	s := []bool{true, false}
	assert.Equal(t,
		[]Node{
			{"code", nil, []any{true}},
			{"code", nil, []any{false}},
		},
		Slice(s, func(i int, b bool) Node {
			return N("code", b)
		}))
}

func TestSliceSeparator(t *testing.T) {
	s := []bool{true, false}
	assert.Equal(t,
		[]Node{
			{"code", nil, []any{true}},
			{"", nil, []any{H("/")}},
			{"code", nil, []any{false}},
		},
		SliceSeparator(s, "/", func(i int, b bool) Node {
			return N("code", b)
		}))
}
