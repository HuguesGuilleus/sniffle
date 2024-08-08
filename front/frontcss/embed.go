package frontcss

import (
	"cmp"
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"slices"
)

//go:embed *.css
var files embed.FS

var Style = func() (all []byte) {
	r := regexp.MustCompile(`(_[\w\d]+)`)
	m := map[string]string{
		"_line":   ".3ex",
		"_spThin": ".7ex",
		"_sp":     "2ex",

		"_colorA":    "#FFCC00", // EU original yellow
		"_colorB":    "#2E98FF",
		"_colorGrey": "lightgrey",
	}

	entries, err := files.ReadDir(".")
	if err != nil {
		panic(err)
	}
	slices.SortFunc(entries, func(a, b fs.DirEntry) int {
		return cmp.Compare(a.Name(), b.Name())
	})
	for _, entry := range entries {
		data, err := fs.ReadFile(files, entry.Name())
		if err != nil {
			panic(err)
		}
		all = append(all, r.ReplaceAllFunc(data, func(s []byte) []byte {
			out, ok := m[string(s)]
			if !ok {
				panic(fmt.Sprintf("Not found %q in file %q", s, entry.Name()))
			}
			return []byte(out)
		})...)
	}

	return
}()
