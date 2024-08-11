"use strict";

((document) => {
	const qsa = (q, f, doc = document) => [...doc.querySelectorAll(q)].map(f),
		lang = document.documentElement.lang,
		dateFormater = new Intl.DateTimeFormat(lang, { dateStyle: "full" }),
		instantFormater = new Intl.DateTimeFormat(lang, {
			dateStyle: "full",
			timeStyle: "long",
		});

	qsa(
		"time",
		(time) =>
			time.innerText =
				(/T/.test(time.dateTime) ? instantFormater : dateFormater)
					.format(new Date(time.dateTime)),
	);

	qsa("input[type=search]", (input) => {
		const items = qsa(
			".si",
			(item) => [
				item,
				qsa(".st", (t) => t.innerText.toLowerCase(), item).join(" "),
			],
		);

		input.hidden = false;
		input.focus();
		input.oninput = () => {
			const queries = input.value.toLowerCase().split(/\s+/);
			items.map(([i, t]) =>
				i.hidden = !queries.every((query) => t.includes(query))
			);
		};
	});
})(document);
