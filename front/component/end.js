"use strict";

((document) => {
	const qsa = (q, f) => document.querySelectorAll(q).forEach(f),
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
				(time.dateTime.includes("T") ? instantFormater : dateFormater)
					.format(new Date(time.dateTime)),
	);
})(document);
