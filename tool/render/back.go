package render

// A static page to redirect to the parent
var Back = []byte(`<!DOCTYPE html><html>` +
	`<meta http-equiv=refresh content='0;URL=..'>` +
	`<body>` +
	`<a id=a href=..>Redirect</a>`)
