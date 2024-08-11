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
	LogoTitle     string   `help:"Logo title and alt attribute"`
	FooterBuild   render.H `help:"In Footer, before build date"`

	ABOUT struct {
		PageTitle       string     `help:"Header title"`
		PageDescription string     `help:"head description"`
		Text            []render.H `help:"page content"`
	} `help:"Website about page"`

	HOME struct {
		Name            string `help:"Full Home text link"`
		PageDescription string `help:"Website home page head description"`
	} `help:"website home page"`

	EU struct {
		Name string `help:"Name of European Union"`
	}

	EU_EC struct {
		Name string `help:"Name of European Commission"`
	}

	EU_EC_ECI struct {
		Name  string `help:"Name of European Citizens' Initiative"`
		INDEX struct {
			Name            string
			PageDescription string
		} `help:"ECI index page"`
		ONE struct {
			Status               render.H
			LastUpdate           render.H
			LinkOfficial         render.H `help:"Official page link"`
			LinkSupport          render.H `help:"Support page link"`
			LinkWebsite          render.H `help:"Organisator website link"`
			H1DescriptionGeneral render.H
			H1DescriptionAnnex   render.H
			H1Treaty             render.H
			H1Signature          render.H
		}
		Status map[string]render.H
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
