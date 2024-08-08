package translate

import (
	_ "embed"
	"encoding/json"
	"sniffle/tool/language"
	"sniffle/tool/render"
)

type Translation struct {
	PageTop       render.H `help:"Page header to indicated that this website is not official"`
	AboutTextLink render.H `help:"About text for link"`

	EU_EC_ICE_INDEX struct{}

	EU_EC_ICE_ONE struct {
		LastUpdate           render.H
		H1DescriptionGeneral render.H
		H1DescriptionAnnex   render.H
		H1Signature          render.H
	}
}

var AllTranslation = map[language.Language]Translation{
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
