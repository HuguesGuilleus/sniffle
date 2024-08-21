package securehtml

import (
	"html"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecure(t *testing.T) {
	assert.Equal(t, template.HTML(``), Secure(``))
	assert.Equal(t, template.HTML(`a <a href="https://google.com">https://google.com</a> b`), Secure(`a https://google.com b`))
	assert.Equal(t, template.HTML(``), Secure(`<hr><script>alert("SECURITY!!!");</script><style>*{display:none!important}</style>`))
	assert.Equal(t, template.HTML(`<p><b>Very</b> <i>Safe</i></p>`), Secure(`<div><p why=42><b>Very</b> <i>Safe</i></p>`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h1>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h2>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h3>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h4>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h5>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h6>Title`))
	assert.Equal(t, template.HTML(`<br>`), Secure(`<br/>`))
	assert.Equal(t, template.HTML(`&gt;&lt;&amp;&#34;&#39;`), Secure(html.EscapeString(`<>&"'`)))
	assert.Equal(t, template.HTML(`<p><b>H</b>ello</p><p>World <b>!</b></p>`), Secure(` <p> <b>H</b>ello </p> <p> World <b> !`))
	assert.Equal(t, template.HTML(`<p>Hello</p>yo<pre>	World</pre>End`), Secure(` <p> Hello </p> yo <pre>	World</pre>	End`))
	// assert.Equal(t, template.HTML(``), Secure(`<p><i></p>`))
	assert.Equal(t, template.HTML(`<a href="https://europa.eu/">Link</a>`), Secure(`<a href="https://europa.eu/">Link`))
	assert.Equal(t, template.HTML(`<a href="http://europa.eu/">Link</a>`), Secure(`<a href="http://europa.eu/">Link`))
	assert.Equal(t, template.HTML(`<del>[file:///home/user/] Link</del>`), Secure(`<a href="file:///home/user/">Link`))
}

func TestText(t *testing.T) {
	assert.Equal(t, ``, Text(``, 10))
	assert.Equal(t, `Title 1`, Text(`<p> <strong>Title   1</strong>No</p>`, 7))
}

func TestURL(t *testing.T) {
	assert.Nil(t, ParseURL(""))
	assert.Nil(t, ParseURL(":x/"))
	assert.Nil(t, ParseURL("file://root"))
	assert.Equal(t, "https://google.com", ParseURL("google.com").String())
}
