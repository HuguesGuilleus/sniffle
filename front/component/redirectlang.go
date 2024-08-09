package component

import (
	"sniffle/tool/language"
	"strings"
)

// Create a redirect page for
func RedirectIndex(langs []language.Language) []byte {
	h := []byte(`<!DOCTYPE html>` +
		`<html>` +
		`<head>` +
		`<meta charset=utf-8>` +
		`<meta name=robots content=noindex>` +
		`</head>` +
		`<body>` +
		`<p>Choose a language:</p>`)

	for _, l := range langs {
		h = append(h, `<a hreflang=`...)
		h = append(h, l.String()...)
		h = append(h, ` href=`...)
		h = append(h, l.String()...)
		h = append(h, `.html>`...)
		h = append(h, l.Human()...)
		h = append(h, `</a><br>`...)
	}

	h = append(h, redirectIndexJs...)

	return h
}

var redirectIndexJs = strings.NewReplacer("\n", "", "\t", "").Replace(`<script>
var m={},
	A=document.querySelectorAll("a"),
	L=navigator.languages,
	i=0,
	a;

for(;i<A.length;i++)
	m[A[i].hreflang]=A[i].href;

for(i=0;i<L.length;i++){
	a=m[L[i].replace(/-\w+/,"")];
	console.log("a:", a, L[i], L[i].replace(/-\w+/,""));
	if(a){location=a;break}
}

</script>`)
