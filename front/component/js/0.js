const INNERTEXT = "innerText",
	qsa = (q, f, doc = document) => [...doc.querySelectorAll(q)].map(f),
	DateTimeFormat = (opt) =>
		new Intl.DateTimeFormat(document.documentElement.lang, {
			dateStyle: "full",
			...opt,
		});

qsa(
	"time",
	(time) =>
		time[INNERTEXT] = (/T/.test(time.dateTime)
			? DateTimeFormat({ timeStyle: "long" })
			: DateTimeFormat({}))
			.format(new Date(time.dateTime)),
);
