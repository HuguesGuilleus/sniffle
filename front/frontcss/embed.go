package frontcss

import (
	"bytes"
	"cmp"
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"slices"

	"github.com/tdewolff/minify/v2/css"
)

//go:embed *.css
var files embed.FS

var Style = func() []byte {
	r := regexp.MustCompile(`(_[\w\d]+)`)
	m := map[string]string{
		"_line":   ".1rem",
		"_spThin": ".5rem",
		"_sp":     "1rem",

		"_back":   "#EEE",
		"_back1":  "#DDD",
		"_color":  "black",
		"_color1": "#222",

		"_colorA": "#2E98FF",
		// "_colorB":    "#FFCC00", // EU original yellow
		// "_colorGrey": "lightgrey",
	}

	entries, err := files.ReadDir(".")
	if err != nil {
		panic(err)
	}
	slices.SortFunc(entries, func(a, b fs.DirEntry) int {
		return cmp.Compare(a.Name(), b.Name())
	})
	buff := bytes.Buffer{}
	for _, entry := range entries {
		data, err := fs.ReadFile(files, entry.Name())
		if err != nil {
			panic(err)
		}
		buff.WriteString(r.ReplaceAllStringFunc(string(data), func(s string) string {
			out, ok := m[s]
			if !ok {
				panic(fmt.Sprintf("Not found %q in file %q", s, entry.Name()))
			}
			return out
		}))
	}

	out := bytes.Buffer{}
	if err := css.Minify(nil, &out, &buff, nil); err != nil {
		panic(err.Error())
	}

	return out.Bytes()
}()
