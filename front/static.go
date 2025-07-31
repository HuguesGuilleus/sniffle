// front contains all assets and tools to render page.
package front

import (
	"embed"

	"github.com/HuguesGuilleus/sniffle/common"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/fronttool"
)

//go:embed favicon.ico
var favicon []byte

//go:embed robots.txt
var robots []byte

var (
	//go:embed frontcss/*
	styleFiles embed.FS
	styleData  = fronttool.CSS(styleFiles, map[string]string{
		"_line1":  ".1rem",
		"_line2":  ".2rem",
		"_spThin": ".5rem",
		"_sp":     "1rem",
		"_spsp":   "2rem",
		"_sp1":    "1.5rem",
		"_sp2":    "2.5rem",

		"_back":   "#EEE",
		"_back1":  "#DDD",
		"_color":  "black",
		"_color1": "#222",
		"_color2": "#555",

		"_colorA":     "#2E98FF",
		"_colorADark": "#006ad0",
		"_colorEdito": "orchid",
	})
	StyleHash, StyleIntegrity = fronttool.FileSum(styleData)
)

func Do(t *tool.Tool) {
	t.WriteFile("favicon.ico", favicon)
	t.WriteFile("robots.txt", append(robots, ("\nSitemap: "+common.Host+"/sitemap.txt\n")...))
	t.WriteFile("style."+StyleHash+".css", styleData)
}
