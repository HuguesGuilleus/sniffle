// language manage UE speaken languages as integer constants.
package language

import (
	"fmt"
	"strings"
)

type Language uint

// Langue constants.
// The order can change.
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
	// Use english, but do not print anchor href to english page, but ./index.html for client redirect to his language.
	AllEnglish

	// The length for a array of language codes.
	Len = AllEnglish + 1
)

// https://en.wikipedia.org/wiki/Languages_of_the_European_Union#Official_EU_languages
var language2iso = [Len]string{
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
	AllEnglish: "en",
}
var language2ISO = [Len]string{
	Bulgarian:  "BG",
	Croatian:   "HR",
	Czech:      "CS",
	Danish:     "DA",
	Dutch:      "NL",
	English:    "EN",
	Estonian:   "ET",
	Finnish:    "FI",
	French:     "FR",
	German:     "DE",
	Greek:      "EL",
	Hungarian:  "HU",
	Irish:      "GA",
	Italian:    "IT",
	Latvian:    "LV",
	Lithuanian: "LT",
	Maltese:    "MT",
	Polish:     "PL",
	Portuguese: "PT",
	Romanian:   "RO",
	Slovak:     "SK",
	Slovene:    "SL",
	Spanish:    "ES",
	Swedish:    "SV",
	AllEnglish: "EN",
}

// Two lower ascii letter of the ISO language code.
// If unknwon return "??".
func (l Language) String() string {
	if l == Invalid || l >= Len {
		return "??"
	}
	return language2iso[l]
}

// Two upper ascii letter of the ISO language code.
// If unknwon return "??".
func (l Language) Upper() string {
	if l == Invalid || l >= Len {
		return "??"
	}
	return language2ISO[l]
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

// UnmarshalText is reverse of Language.String().
// It take two letter (case insensitive) and convert is language.
func (l *Language) UnmarshalText(data []byte) error {
	s := string(data)
	*l = iso2language[strings.ToLower(s)]
	if *l == Invalid {
		return fmt.Errorf("unknwon langage %q", s)
	}
	return nil
}

var language2human = [Len]string{
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
	AllEnglish: "English",
}

// The name of the langue in this langue.
func (l Language) Human() string {
	if l == Invalid || l >= Len {
		return "_LANGUAGE_"
	}
	return language2human[l]
}

// Create path.
// If l == [AllEnglish], return without change.
// Is this case, the basePath must ends with a '/', else panic.
func (l Language) Path(basePath string) string {
	if l == AllEnglish {
		if !strings.HasSuffix(basePath, "/") {
			panic(fmt.Sprintf("Language.Path(%q) do not end with a '/'", basePath))
		}
		return basePath
	}
	return basePath + l.String() + ".html"
}
