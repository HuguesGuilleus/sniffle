document.querySelectorAll("a").forEach((a) =>
	a.onclick = (_) => {
		const doc = window.open().document;
		doc.title = a.dataset.title;
		doc.body.innerText = atob(a.dataset.b64);
		doc.body.style.fontFamily = "monospace";
		doc.body.style.fontSize = "x-large";
		doc.body.style.background = "#EED";
		doc.body.style.background = "#EED";
		doc.body.style.wordBreak = "break-word";
		doc.body.style.whiteSpace = "pre-line";
	}
);

document.querySelectorAll("code").forEach((c) =>
	c.onclick = (_) =>
		navigator.clipboard.writeText(c.dataset.id) | alert("copied!")
);

const lis = [...document.querySelectorAll("li")].map(
	(li) => [li, li.innerText.toLowerCase()],
);
s.oninput = (_, p = s.value.toLowerCase().split(/\s+/)) =>
	lis.forEach(([li, t]) => li.hidden = !p.every((p) => t.includes(p)));
