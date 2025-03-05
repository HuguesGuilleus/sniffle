const tocItems = [], // []->[tocItem, [correspondingElements...]]
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

let currentTocElement1 = [],
	currentTocElement2 = [],
	tagName;

qsa("*", (element) => {
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
