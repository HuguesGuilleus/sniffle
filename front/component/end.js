"use strict";

((document) => {
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

	qsa("input[type=search]", (input) => {
		const items = qsa(
			".si",
			(item) => [
				item,
				qsa(".st", (t) => t[INNERTEXT].toLowerCase(), item).join(" "),
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

	qsa(".wt", (_) => {
		let tocAutoId = 1,
			currentTocElement1 = [],
			currentTocElement2 = [];
		//  []->[tocItem, [correspondingElements...]]
		const tocItems = [],
			tocItemsPush = (
				element,
				level,
				array,
				tocItem = document.createElement("a"),
			) => {
				tocItem.className = "wi wi" + level;
				tocItem.href = "#" + (element.id ||= tocAutoId++);
				tocItem[INNERTEXT] = element[INNERTEXT];
				toc.append(tocItem);
				tocItems.push([tocItem, array]);
			},
			visibleElement = new Map(),
			observer = new IntersectionObserver((entries) => {
				for (const entry of entries) {
					visibleElement.set(entry.target, entry.isIntersecting);
				}
				for (const [tocItem, elements] of tocItems) {
					tocItem.dataset.w = elements.some((element) =>
						visibleElement.get(element)
					);
				}
			});

		toc[INNERTEXT] = "";
		qsa(".wc>*", (element) => {
			if (element.tagName == "H1") {
				tocItemsPush(element, 1, currentTocElement1 = []);
				currentTocElement2 = [];
			} else if (element.tagName == "H2") {
				tocItemsPush(element, 2, currentTocElement2 = []);
			}
			currentTocElement1.push(element);
			currentTocElement2.push(element);
			observer.observe(element);
		});
	});
})(document);
