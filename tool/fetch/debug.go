package fetch

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sniffle/tool/render"
	"strconv"
	"time"
)

const (
	DebugKeepIgnore = iota
	DebugKeepIndex
	DebugKeepData
)

type deebugfile struct {
	*Meta
	txt string
}

// Debug create cacheRoot/index.html with a index of all cache request.
// The keeper return a const like [DebugKeepIgnore].
func Debug(cacheRoot string, keeper func(host string) int) error {
	index := []*deebugfile{}

	paths, err := filepath.Glob(filepath.Join(cacheRoot, filepath.FromSlash("/*/*/*.http")))
	if err != nil {
		return err
	}
	for _, p := range paths {
		meta, err := debugRead(p, keeper)
		if err != nil {
			return err
		}
		if meta != nil {
			index = append(index, meta)
		}
	}

	slices.SortFunc(index, func(a, b *deebugfile) int {
		return cmp.Or(
			cmp.Compare(a.RawURL, b.RawURL),
			cmp.Compare(a.ID(), b.ID()),
		)
	})

	os.WriteFile(
		filepath.Join(cacheRoot, "index.html"),
		debugRender(index),
		0o664,
	)

	return nil
}

func debugRead(path string, keeper func(host string) int) (*deebugfile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	meta, err := ReadOnlyMeta(f)
	if err != nil {
		return nil, fmt.Errorf("read cache %q: %w", path, err)
	}

	switch keeper(meta.URL.Host) {
	case DebugKeepIgnore:
		return nil, nil
	case DebugKeepIndex:
		return &deebugfile{
			Meta: meta,
		}, nil
	}

	begin := [4096]byte{}
	n, _ := f.Read(begin[:])

	return &deebugfile{
		Meta: meta,
		txt:  debugPrint(meta, begin[:n]),
	}, nil
}

func debugPrint(meta *Meta, begin []byte) string {
	buff := bytes.Buffer{}

	buff.WriteString(meta.Path())
	buff.WriteString("\n")
	buff.WriteString(meta.Time.UTC().Format(time.DateTime))
	buff.WriteString("\n\n")

	buff.WriteString(meta.Method)
	buff.WriteString(" ")
	buff.WriteString(meta.URL.String())
	buff.WriteString("\n")
	meta.RequestHeader.Write(&buff)
	buff.WriteString("\n")

	if len(meta.RequestBody) != 0 {
		buff.Write(meta.RequestBody)
		buff.WriteString("\n\n")
	}

	buff.WriteString(strconv.Itoa(meta.Status))
	buff.WriteString("\n")
	meta.ResponseHeader.Write(&buff)
	buff.WriteString("\n")

	buff.Write(begin)

	return buff.String()
}

//go:embed debug.css
var debugcss render.H

//go:embed debug.js
var debugjs render.H

func debugRender(index []*deebugfile) []byte {
	return render.Merge(render.N("html",
		render.N("head",
			render.H(`<meta charset=utf-8>`),
			render.H(`<meta name=viewport content="width=device-width,initial-scale=1">`),
			render.N("title", "Index of HTTP cache"),
			render.N("style", debugcss),
		),
		render.N("body",
			render.H(`<input id=s type=search placeholder="Search ...">`),
			render.N("ul", render.S(index, "", func(f *deebugfile) render.Node {
				return render.N("li",
					render.Na("code", "data-id", f.Path()).N(f.ID()[:7]),
					" ",
					render.If(len(f.txt) > 0, func() render.Node {
						return render.N("",
							render.Na("a", "data-b64", base64.StdEncoding.EncodeToString([]byte(f.txt))).
								A("data-title", f.URL.Host+"/"+f.ID()[:7]).
								N("~>"),
							" ",
						)
					}),
					render.IfElseS(f.txt == "", "-- ", ""),
					render.IfS(f.Status/100 == 2, render.N("span.s2", "s", f.Status)),
					render.IfS(f.Status/100 == 3, render.N("span.s3", "s", f.Status)),
					render.IfS(f.Status/100 == 4, render.N("span.s4", "s", f.Status)),
					render.IfS(f.Status/100 == 5, render.N("span.s5", "s", f.Status)),
					" ",
					render.N("span.m", f.Method),
					" ",
					f.URL,
					render.N(""),
				)
			})),
			render.N("script", debugjs),
		),
	))
}
