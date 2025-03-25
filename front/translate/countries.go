package translate

import (
	"cmp"
	"slices"
	"sniffle/common/country"
	"sniffle/common/language"
)

// Get countries sorted in language l.
func Countries(l language.Language) []country.Country {
	countries := []country.Country{
		country.Austria,
		country.Belgium,
		country.Bulgaria,
		country.Croatia,
		country.Cyprus,
		country.Czechia,
		country.Denmark,
		country.Estonia,
		country.Finland,
		country.France,
		country.Germany,
		country.Greece,
		country.Hungary,
		country.Ireland,
		country.Italy,
		country.Latvia,
		country.Lithuania,
		country.Luxembourg,
		country.Malta,
		country.Netherlands,
		country.Poland,
		country.Portugal,
		country.Romania,
		country.Slovakia,
		country.Slovenia,
		country.Spain,
		country.Sweden,
		country.UnitedKingdom,
	}
	tr := T[l]
	slices.SortFunc(countries, func(a, b country.Country) int {
		return cmp.Compare(tr.Country[a], tr.Country[b])
	})
	return countries
}
