// DOM architecture:
// input[type=search]
// .sg [ // search group: if any search item is display, hide the group
//     .si [ // if any queries match searcg target, hide this item.
//         .st // search target: use this string to math the queries.
//     ]
// ]
// .sg []

// Type: [][searchGroup, [][searchItem, searchTarget:string]]
const searchGroupArray = qsa(
		".sg",
		(searchGroup) => [
			searchGroup,
			qsa(".si", (searchItem) => [
				searchItem,
				qsa(
					".st",
					(searchTarget) => searchTarget[INNERTEXT].toLowerCase(),
					searchItem,
				).join(" "),
			], searchGroup),
		],
	),
	_search = qsa("[type=search]", (input) => {
		input[HIDDEN] = false;
		input.focus();
		input.oninput = (
			_,
			queries = input.value.toLowerCase().split(/\s+/),
		) => searchGroupArray.map(([searchGroup, items]) =>
			searchGroup[HIDDEN] = items.map(([searchItem, searchTarget]) =>
				searchItem[HIDDEN] = queries.some((query) =>
					!searchTarget.includes(query)
				)
			).every((b) => b)
		);
	});
