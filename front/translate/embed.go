package translate

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"sniffle/tool/language"
)

type Translation struct {
	PageTop       template.HTML `help:"Page header to indicated that this website is not official"`
	AboutTextLink template.HTML `help:"About text for link"`
}

var AllTranslation = map[language.Langage]Translation{
	language.English: load(fileEn),
	language.French:  load(fileFR),
}

var (
	//go:embed translate.en.json
	fileEn []byte
	//go:embed translate.fr.json
	fileFR []byte
)

func load(data []byte) (t Translation) {
	if err := json.Unmarshal(data, &t); err != nil {
		panic(err)
	}
	return
}
