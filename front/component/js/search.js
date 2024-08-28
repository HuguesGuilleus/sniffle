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
	input.oninput = (_) => {
		const queries = input.value.toLowerCase().split(/\s+/);
		items.map(([i, t]) =>
			i.hidden = !queries.every((query) => t.includes(query))
		);
	};
});
