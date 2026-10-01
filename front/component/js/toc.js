const tocItems = [],
	tocItemsPush = (
		header,
		level,
		array,
		tocItem = document.createElement("a"),
	) => {
		tocItem.className = "wi wi" + level;
		tocItem.href = "#" + (header.id ||= header[INNERTEXT]);
		tocItem[INNERTEXT] = header[INNERTEXT];
		toc.append(tocItem);
		tocItems.push([tocItem, header, array]);
	},
	visibleElement = new Map(),
	observer = new IntersectionObserver((entries) =>
		entries.map((entry) =>
			visibleElement.set(entry.target, entry.isIntersecting)
		) |
		tocItems.map(([tocItem, header, elements]) => {
			tocItem[HIDDEN] = !header.offsetParent;
			tocItem.dataset.w = elements.some((element) =>
				visibleElement.get(element)
			);
		})
	);

let currentTocElement2 = [],
	currentTocElement3 = [],
	tagName;

qsa("*", (element) => {
	tagName = element.tagName;
	if (tagName == "H2") {
		tocItemsPush(element, 2, currentTocElement2 = []);
		currentTocElement3 = [];
	} else if (tagName == "H3") {
		tocItemsPush(element, 3, currentTocElement3 = []);
	}
	currentTocElement2.push(element);
	currentTocElement3.push(element);
	observer.observe(element);
});
