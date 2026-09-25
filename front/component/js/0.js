const INNERTEXT = "innerText",
	HIDDEN = "hidden",
	qsa = (q, f, doc = document) => [...doc.querySelectorAll(q)].map(f);
