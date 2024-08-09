package translate

import (
	_ "embed"
	"encoding/json"
	"sniffle/tool/language"
	"sniffle/tool/render"
)

type Translation struct {
	PageTop       render.H `help:"Page header to indicated that this website is not official"`
	InDev         render.H `help:"This page is actualy on development"`
	AboutTextLink render.H `help:"About text for link"`
	HomeTitle     string   `help:"Home page link title"`
	HomeName      render.H `help:"Home page link name"`
	FooterBuild   render.H `help:"In Footer, before build date"`

	ABOUT struct {
		PageTitle       string     `help:"Header title"`
		PageDescription string     `help:"head description"`
		Text            []render.H `help:"page content"`
	} `help:"Website about page"`

	EU struct {
		Name string `help:"Name of European Union"`
	}

	EU_EC struct {
		Name string `help:"Name of European Commission"`
	}

	EU_EC_ICE struct {
		Name  string `help:"Name of European Citizens' Initiative"`
		INDEX struct{}
		ONE   struct {
			LastUpdate           render.H
			H1DescriptionGeneral render.H
			H1DescriptionAnnex   render.H
			H1Treaty             render.H
			H1Signature          render.H
		}
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
