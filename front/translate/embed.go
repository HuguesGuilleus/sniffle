// translate contains texts in multiple languages.
package translate

import (
	_ "embed"
	"encoding/json"
	"sniffle/common/country"
	"sniffle/common/language"
	"sniffle/tool/render"
)

type Translation struct {
	Global struct {
		Presentation render.H
	}

	AboutTextLink render.H `help:"About text for link"`
	FooterBuild   render.H `help:"In Footer, before build date"`
	InDev         render.H `help:"This page is actualy on development"`
	LinkOfficial  render.H `help:"Official page link"`
	LogoTitle     string   `help:"Logo title and alt attribute"`
	PageTop       render.H `help:"Page header to indicated that this website is not official"`
	SchemaLink    render.H `help:"Schema link"`
	SearchInside  string
	HELP          render.H
	Byte          render.H
	Source        render.H

	Country [country.Len]render.H  `json:"-"`
	Langage [language.Len]render.H `json:"-"`

	ABOUT struct {
		PageTitle       string     `help:"Header title"`
		PageDescription string     `help:"head description"`
		Intro           []render.H `help:"page content"`
		Mail            render.H
		Host            render.H
	} `help:"Website about page"`

	HOME struct {
		Name            string `help:"Full Home text link"`
		PageDescription string `help:"Website home page head description"`
		EU              string `help:"European union section"`
		About           string `help:"About section"`
		Release         string `help:"Version release (in english)"`
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
			Help            render.H
		} `help:"ECI index page"`
		ONE struct {
			Status                     render.H
			Categorie                  render.H
			LastUpdate                 render.H
			DescriptionOriginalLangage render.H
			LinkSignature              render.H `help:"Link to signature page"`
			LinkFollowUp               render.H
			LinkWebsite                render.H `help:"Link to organisator"`
			H1Description              render.H
			H1DescriptionAnnex         render.H
			AnnexDocument              render.H
			DraftLegal                 render.H
			H1Treaty                   render.H
			H1Timeline                 render.H
			Registration               render.H
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

			AnswerKind struct {
				Annex        render.H
				Response     render.H
				PressRelease render.H
			}

			Member struct {
				H1      render.H
				Type    map[string]render.H
				Start   render.H
				End     render.H
				Country render.H
				Privacy render.H
			}

			Funding struct {
				Name               render.H
				Total              render.H
				Sponsor            render.H
				Amount             render.H
				Date               render.H
				PrivateSponsor     render.H
				PrivateSponsorHelp string
				CaptionDate        render.H
				CaptionAmount      render.H
				Document           render.H
			}
		}
		Status        map[string]render.H
		Categorie     map[string]render.H
		ThresholdRule map[string]render.H
	}
}

var Langs = []language.Language{
	language.English,
	language.French,
}

var T = [...]Translation{
	language.English: load(fileEn),
	language.French:  load(fileFR),

	language.AllEnglish: load(fileEn),
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
