qsa(".wt", (_) => {
	let tocAutoId = 1,
		currentTocElement1 = [],
		currentTocElement2 = [],
		tagName;

	const tocItems = [], // []->[tocItem, [correspondingElements...]]
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
	qsa(".wc *", (element) => {
		tagName = element.tagName;
		if (tagName == "H1") {
			tocItemsPush(element, 1, currentTocElement1 = []);
			currentTocElement2 = [];
		} else if (tagName == "H2") {
			tocItemsPush(element, 2, currentTocElement2 = []);
		}
		currentTocElement1.push(element);
		currentTocElement2.push(element);
		observer.observe(element);
	});
});
