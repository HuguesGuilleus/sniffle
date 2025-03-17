// country manage UE Country as integer constants.
package country

import (
	"fmt"
)

type Country uint

// The List of Country.
// The code can change, any no order
// Can contain non UE country.
const (
	Invalid Country = iota
	Austria
	Belgium
	Bulgaria
	Croatia
	Cyprus
	Czechia
	Denmark
	Estonia
	Finland
	France
	Germany
	Greece
	Hungary
	Ireland
	Italy
	Latvia
	Lithuania
	Luxembourg
	Malta
	Netherlands
	Poland
	Portugal
	Romania
	Slovakia
	Slovenia
	Spain
	Sweden
	UnitedKingdom

	// The length for a array of coutrie codes.
	Len = UnitedKingdom + 1
)

func (c Country) IsZero() bool {
	return c < Austria || UnitedKingdom < c
}
func (c Country) NotZero() bool {
	return Austria <= c && c <= UnitedKingdom
}

func (c *Country) UnmarshalText(data []byte) error {
	s := string(data)
	*c = fromJSON[s]
	if *c == Invalid {
		return fmt.Errorf("unknwon country: %q", s)
	}

	return nil
}

func (c Country) String() string {
	s := country2iso[c]
	if s == "" {
		return "??"
	}
	return s
}

var fromJSON = map[string]Country{
	"AT": Austria,
	"BE": Belgium,
	"BG": Bulgaria,
	"CY": Cyprus,
	"CZ": Czechia,
	"DE": Germany,
	"DK": Denmark,
	"EE": Estonia,
	"ES": Spain,
	"FI": Finland,
	"FR": France,
	"GR": Greece,
	"HR": Croatia,
	"HU": Hungary,
	"IE": Ireland,
	"IT": Italy,
	"LT": Lithuania,
	"LU": Luxembourg,
	"LV": Latvia,
	"MT": Malta,
	"NL": Netherlands,
	"PL": Poland,
	"PT": Portugal,
	"RO": Romania,
	"SE": Sweden,
	"SI": Slovenia,
	"SK": Slovakia,
	"GB": UnitedKingdom,

	"at": Austria,
	"be": Belgium,
	"bg": Bulgaria,
	"cy": Cyprus,
	"cz": Czechia,
	"de": Germany,
	"dk": Denmark,
	"ee": Estonia,
	"es": Spain,
	"fi": Finland,
	"fr": France,
	"gr": Greece,
	"hr": Croatia,
	"hu": Hungary,
	"ie": Ireland,
	"it": Italy,
	"lt": Lithuania,
	"lu": Luxembourg,
	"lv": Latvia,
	"mt": Malta,
	"nl": Netherlands,
	"pl": Poland,
	"pt": Portugal,
	"ro": Romania,
	"se": Sweden,
	"si": Slovenia,
	"sk": Slovakia,
	"gb": UnitedKingdom,

	"Austria":       Austria,
	"Belgium":       Belgium,
	"Bulgaria":      Bulgaria,
	"Cyprus":        Cyprus,
	"Czechia":       Czechia,
	"Germany":       Germany,
	"Denmark":       Denmark,
	"Estonia":       Estonia,
	"Spain":         Spain,
	"Finland":       Finland,
	"France":        France,
	"Greece":        Greece,
	"Croatia":       Croatia,
	"Hungary":       Hungary,
	"Ireland":       Ireland,
	"Italy":         Italy,
	"Lithuania":     Lithuania,
	"Luxembourg":    Luxembourg,
	"Latvia":        Latvia,
	"Malta":         Malta,
	"Netherlands":   Netherlands,
	"Poland":        Poland,
	"Portugal":      Portugal,
	"Romania":       Romania,
	"Sweden":        Sweden,
	"Slovenia":      Slovenia,
	"Slovakia":      Slovakia,
	"UnitedKingdom": UnitedKingdom,
}

var country2iso = map[Country]string{
	Austria:       "AT",
	Belgium:       "BE",
	Bulgaria:      "BG",
	Cyprus:        "CY",
	Czechia:       "CZ",
	Germany:       "DE",
	Denmark:       "DK",
	Estonia:       "EE",
	Spain:         "ES",
	Finland:       "FI",
	France:        "FR",
	Greece:        "GR",
	Croatia:       "HR",
	Hungary:       "HU",
	Ireland:       "IE",
	Italy:         "IT",
	Lithuania:     "LT",
	Luxembourg:    "LU",
	Latvia:        "LV",
	Malta:         "MT",
	Netherlands:   "NL",
	Poland:        "PL",
	Portugal:      "PT",
	Romania:       "RO",
	Sweden:        "SE",
	Slovenia:      "SI",
	Slovakia:      "SK",
	UnitedKingdom: "GB",
}
