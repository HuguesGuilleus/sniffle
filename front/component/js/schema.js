qsa(
	"a.block",
	(a) =>
		a.onclick = (event) =>
			event.preventDefault() | navigator.clipboard.writeText(a.href),
);
