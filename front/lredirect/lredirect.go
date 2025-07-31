// lredirect generate pages for index.html or by default when the language isn't available.
package lredirect

import (
	"bytes"
	"html"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/translate"
)

// A page with all available translate languages. No official link.
var All = Page("", translate.Langs)

// Page generates a basic html page with the offical page link if any, and links for all languages.
// The official page link is escaped.
func Page(official string, langs []language.Language) []byte {
	buff := bytes.Buffer{}

	buff.WriteString(`` +
		`<!DOCTYPE html>` +
		`<html>` +
		`<head>` +
		`<meta charset=utf-8>` +
		`<meta name=robots content=noindex>`)

	if len(langs) != 0 {
		buff.WriteString(`` +
			`<script>` +
			`for(a of navigator.languages)` +
			`if(` +
			`/`)
		for i, l := range langs {
			if i != 0 {
				buff.WriteByte('|')
			}
			buff.WriteString(l.String())
		}
		buff.WriteString(`` +
			`/.test(a=a.split("-")[0])` +
			`){location=a+".html";break}` +
			`</script>`)
	}

	buff.WriteString(`` +
		`<meta name=viewport content="width=device-width,initial-scale=1">` +
		`<style>` +
		`a{display:inline-block;padding:2ex;font-size:xx-large}` +
		`</style>` +
		`</head>` +
		`<body>`)

	if official != "" {
		official = html.EscapeString(official)
		buff.WriteString(`<a href="`)
		buff.WriteString(official)
		buff.WriteString(`">Official Page: `)
		buff.WriteString(official)
		buff.WriteString(`</a><hr>`)
	}

	for _, l := range langs {
		buff.WriteString(`<a href=`)
		buff.WriteString(l.String())
		buff.WriteString(`.html>`)
		buff.WriteString(l.Human())
		buff.WriteString(`</a>`)
	}

	return buff.Bytes()
}
