package render

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	h := Merge(No("body.cl1.cl2.cl3",
		A("hidden", "").
			A("data-toescape", "<must' escape\">").
			A("data-normal", "yolo/"),
		"Hello ",
		N("i.cl", "World"),
		H("&#33;"),
		"\n",
		true,
		uint(42),
		-1,
		boolHTML(true),
	))

	assert.Equal(t, `<!DOCTYPE html>`+
		`<body class="cl1 cl2 cl3"`+
		` hidden`+
		` data-toescape="<must' escape&#34;>"`+
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
		`</body>`,
		string(h))
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

type boolHTML bool

func (b boolHTML) HTML() H {
	if b {
		return "TRUE"
	} else {
		return "FALSE"
	}
}
