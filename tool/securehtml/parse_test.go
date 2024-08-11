package securehtml

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecure(t *testing.T) {
	assert.Equal(t, template.HTML(``), Secure(``))
	assert.Equal(t, template.HTML(``), Secure(`<br><script>alert("SECURITY!!!");</script><style>*{display:none!important}</style>`))
	assert.Equal(t, template.HTML(`<p><b>Very</b> <i>Safe</i></p>`), Secure(`<div><p why=42><b>Very</b> <i>Safe</i></p>`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h1>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h2>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h3>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h4>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h5>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), Secure(`<h6>Title`))
	assert.Equal(t, template.HTML(`<a href="https://europa.eu/">Link</a>`), Secure(`<a href="https://europa.eu/">Link`))
	assert.Equal(t, template.HTML(`<a href="http://europa.eu/">Link</a>`), Secure(`<a href="http://europa.eu/">Link`))
	assert.Equal(t, template.HTML(`<del>[file:///home/user/] Link</del>`), Secure(`<a href="file:///home/user/">Link`))
}

func TestText(t *testing.T) {
	assert.Equal(t, ``, Text(``, 10))
	assert.Equal(t, `Title 1`, Text(`<p> <strong>Title   1</strong>No</p>`, 7))
}
