package component

import (
	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/translate"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

type Legal struct {
	Prefix string
	Num    string
	CELEX  string
}

func (legal Legal) Render(l language.Language) render.Node {
	return render.Na("a",
		"href", "https://eur-lex.europa.eu/legal-content/"+l.Upper()+"/TXT/?uri=CELEX:"+legal.CELEX,
	).N(
		translate.T[l].Legal[legal.Prefix][0],
		legal.Num,
		translate.T[l].Legal[legal.Prefix][1],
	)
}
