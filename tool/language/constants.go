package language

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Langage uint

const (
	Invalid Langage = iota
	Bulgarian
	Croatian
	Czech
	Danish
	Dutch
	English
	Estonian
	Finnish
	French
	German
	Greek
	Hungarian
	Irish
	Italian
	Latvian
	Lithuanian
	Maltese
	Polish
	Portuguese
	Romanian
	Slovak
	Slovene
	Spanish
	Swedish
)

// https://en.wikipedia.org/wiki/Languages_of_the_European_Union#Official_EU_languages
var language2iso = map[Langage]string{
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
var iso2langage = map[string]Langage{
	"bg": Bulgarian,
	"hr": Croatian,
	"cs": Czech,
	"da": Danish,
	"nl": Dutch,
	"en": English,
	"et": Estonian,
	"fi": Finnish,
	"fr": French,
	"de": German,
	"el": Greek,
	"hu": Hungarian,
	"ga": Irish,
	"it": Italian,
	"lv": Latvian,
	"lt": Lithuanian,
	"mt": Maltese,
	"pl": Polish,
	"pt": Portuguese,
	"ro": Romanian,
	"sk": Slovak,
	"sl": Slovene,
	"es": Spanish,
	"sv": Swedish,
}

func (l *Langage) UnmarshalJSON(data []byte) error {
	s := ""
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("unmarshal langage code string: %w", err)
	}

	*l = iso2langage[strings.ToLower(s)]
	if *l == Invalid {
		return fmt.Errorf("unknwon langage %q", s)
	}

	return nil
}

func (l Langage) String() string {
	s := language2iso[l]
	if s == "" {
		return "??"
	}
	return s
}
