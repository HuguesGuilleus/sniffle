package render

// A static page to redirect to the parent
var Back = []byte(`<!DOCTYPE html><html>` +
	`<script>` +
	`location=".."` +
	`</script>` +
	`<body>` +
	`<a id=a href=..>Redirect</a>`)
