package language

import (
	"fmt"
	"strings"
)

type Language uint

// Langue constants.
// The order can change, and can be undense.
// Invalid is always == 0
const (
	Invalid Language = iota
	// Sorted by name in own language
	Bulgarian
	Spanish
	Czech
	Danish
	German
	Estonian
	Greek
	English
	French
	Irish
	Croatian
	Italian
	Latvian
	Lithuanian
	Hungarian
	Maltese
	Dutch
	Polish
	Portuguese
	Romanian
	Slovak
	Slovene
	Finnish
	Swedish

	// The length for a array of language codes.
	Len = Swedish + 1
)

// https://en.wikipedia.org/wiki/Languages_of_the_European_Union#Official_EU_languages
var language2iso = map[Language]string{
	Bulgarian:  "bg",
	Croatian:   "hr",
	Czech:      "cs",
	Danish:     "da",
	Dutch:      "nl",
	English:    "en",
	Estonian:   "et",
	Finnish:    "fi",
	French:     "fr",
	German:     "de",
	Greek:      "el",
	Hungarian:  "hu",
	Irish:      "ga",
	Italian:    "it",
	Latvian:    "lv",
	Lithuanian: "lt",
	Maltese:    "mt",
	Polish:     "pl",
	Portuguese: "pt",
	Romanian:   "ro",
	Slovak:     "sk",
	Slovene:    "sl",
	Spanish:    "es",
	Swedish:    "sv",
}

// Two ascii letter of the ISO language code.
// If unknwo return "??"
func (l Language) String() string {
	s := language2iso[l]
	if s == "" {
		return "??"
	}
	return s
}

var iso2language = map[string]Language{
	"bg": Bulgarian,
	"cs": Czech,
	"da": Danish,
	"de": German,
	"el": Greek,
	"en": English,
	"es": Spanish,
	"et": Estonian,
	"fi": Finnish,
	"fr": French,
	"ga": Irish,
	"hr": Croatian,
	"hu": Hungarian,
	"it": Italian,
	"lt": Lithuanian,
	"lv": Latvian,
	"mt": Maltese,
	"nl": Dutch,
	"pl": Polish,
	"pt": Portuguese,
	"ro": Romanian,
	"sk": Slovak,
	"sl": Slovene,
	"sv": Swedish,
}

func (l *Language) UnmarshalText(data []byte) error {
	s := string(data)
	*l = iso2language[strings.ToLower(s)]
	if *l == Invalid {
		return fmt.Errorf("unknwon langage %q", s)
	}
	return nil
}

var language2human = map[Language]string{
	Bulgarian:  "Български",
	Spanish:    "Español",
	Czech:      "Čeština",
	Danish:     "Dansk",
	German:     "Deutsch",
	Estonian:   "Eesti",
	Greek:      "Ελληνικά",
	English:    "English",
	French:     "Français",
	Irish:      "Gaeilge",
	Croatian:   "Hrvatski",
	Italian:    "Italiano",
	Latvian:    "Latviešu",
	Lithuanian: "Lietuvių",
	Hungarian:  "Magyar",
	Maltese:    "Malti",
	Dutch:      "Nederlands",
	Polish:     "Polski",
	Portuguese: "Português",
	Romanian:   "Română",
	Slovak:     "Slovenčina",
	Slovene:    "Slovenščina",
	Finnish:    "Suomi",
	Swedish:    "Svenska",
}

// The name of the langue in this langue.
func (l Language) Human() string {
	s := language2human[l]
	if s == "" {
		return "_LANGUAGE_"
	}
	return s
}
