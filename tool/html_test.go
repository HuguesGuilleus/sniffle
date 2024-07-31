package tool

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHtml(t *testing.T) {
	assert.Equal(t, template.HTML(``), SecureHTML(``))
	assert.Equal(t, template.HTML(``), SecureHTML(`<br><script>alert("SECURITY!!!");</script><style>*{display:none!important}</style>`))
	assert.Equal(t, template.HTML(`<p><b>Very</b> <i>Safe</i></p>`), SecureHTML(`<div><p why=42><b>Very</b> <i>Safe</i></p>`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), SecureHTML(`<h1>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), SecureHTML(`<h2>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), SecureHTML(`<h3>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), SecureHTML(`<h4>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), SecureHTML(`<h5>Title`))
	assert.Equal(t, template.HTML(`<p><strong>Title</strong></p>`), SecureHTML(`<h6>Title`))
	assert.Equal(t, template.HTML(`<a href="https://europa.eu/">Link</a>`), SecureHTML(`<a href="https://europa.eu/">Link`))
	assert.Equal(t, template.HTML(`<a href="http://europa.eu/">Link</a>`), SecureHTML(`<a href="http://europa.eu/">Link`))
	assert.Equal(t, template.HTML(`<del>[file:///home/user/] Link</del>`), SecureHTML(`<a href="file:///home/user/">Link`))
}
