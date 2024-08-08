package render

// A static page to redirect to the parent
const Back = `<!DOCTYPE html><html>` +
	`<script>` +
	`location=".."` +
	`</script>` +
	`<body>` +
	`<a id=a href=..>Redirect</a>`
