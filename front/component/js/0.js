const INNERTEXT = "innerText",
	HIDDEN = "hidden",
	qsa = (q, f, doc = document) => [...doc.querySelectorAll(q)].map(f),
	DateTimeFormat = (opt) =>
		new Intl.DateTimeFormat(document.documentElement.lang, {
			dateStyle: "full",
			...opt,
		}),
	_time = qsa(
		"time",
		(time) =>
			time[INNERTEXT] = (
				/ /.test(time[INNERTEXT])
					? DateTimeFormat({ timeStyle: "long" })
					: /_/.test(time[INNERTEXT])
					? DateTimeFormat({ dateStyle: "short" })
					: DateTimeFormat(0)
			).format(new Date(time.dateTime)),
	);
