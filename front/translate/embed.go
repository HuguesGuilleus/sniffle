package translate

import (
	_ "embed"
	"encoding/json"
	"sniffle/tool/country"
	"sniffle/tool/language"
	"sniffle/tool/render"
)

type Translation struct {
	AboutTextLink render.H `help:"About text for link"`
	FooterBuild   render.H `help:"In Footer, before build date"`
	InDev         render.H `help:"This page is actualy on development"`
	LinkOfficial  render.H `help:"Official page link"`
	LogoTitle     string   `help:"Logo title and alt attribute"`
	PageTop       render.H `help:"Page header to indicated that this website is not official"`
	SearchInside  string
	HELP          render.H
	Byte          render.H
	Source        render.H

	Country [country.Len]render.H  `json:"-"`
	Langage [language.Len]render.H `json:"-"`

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
			Status                     render.H
			Categorie                  render.H
			LastUpdate                 render.H
			DescriptionOriginalLangage render.H
			LinkSupport                render.H `help:"Support page link"`
			LinkWebsite                render.H `help:"Organisator website link"`
			H1Description              render.H
			H1DescriptionAnnex         render.H
			AnnexDocument              render.H
			H1Treaty                   render.H
			H1Timeline                 render.H
			CollectionEarlyClosure     render.H

			H1Signature           render.H
			SignatureSum          render.H
			ValidatedSignature    render.H
			PaperSignaturesUpdate render.H
			CountryOverThreshold  render.H
			Country               render.H
			Signature             render.H
			Threshold             render.H
			OverThreshold         string
		}
		Status        map[string]render.H
		Categorie     map[string]render.H
		ThresholdRule map[string]render.H
	}
}

var AllTranslation = [...]Translation{
	language.English: load(fileEn),
	language.French:  load(fileFR),
}

var (
	//go:embed translate.en.json
	fileEn []byte
	//go:embed translate.fr.json
	fileFR []byte
)

func load(data []byte) Translation {
	dto := struct {
		Translation
		Country map[country.Country]render.H   `json:"$Country"`
		Langage map[language.Language]render.H `json:"$Langage"`
	}{}
	if err := json.Unmarshal(data, &dto); err != nil {
		panic(err)
	}

	for c, html := range dto.Country {
		dto.Translation.Country[c] = html
	}

	for l, html := range dto.Langage {
		dto.Translation.Langage[l] = html
	}

	return dto.Translation
}
