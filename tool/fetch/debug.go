package fetch

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/base64"
	"fmt"
	"slices"
	"sniffle/tool/render"
	"sniffle/tool/writefs"
	"strconv"
	"time"
)

const (
	DebugKeepIgnore = iota
	DebugKeepIndex
	DebugKeepData
)

// Debug create cacheRoot/index.html with a index of all cache request.
// The keeper return a const like [DebugKeepIgnore].
func Debug(fsys writefs.CompleteFS, keeper func(m *Meta) int) error {
	paths, err := indexHTTPFiles(fsys)
	if err != nil {
		return err
	}

	index := []*Meta{}
	for _, p := range paths {
		meta, err := debugRead(fsys, p, keeper)
		if err != nil {
			return err
		} else if meta != nil {
			index = append(index, meta)
		}
	}

	slices.SortFunc(index, func(a, b *Meta) int {
		return cmp.Or(
			cmp.Compare(a.RawURL, b.RawURL),
			cmp.Compare(a.ID(), b.ID()),
		)
	})

	return writefs.WriteFile(fsys, "index.html", debugRender(index))
}

func debugRead(fsys writefs.Opener, path string, keeper func(*Meta) int) (*Meta, error) {
	f, err := fsys.Open(path)
	if err != nil {
		return nil, fmt.Errorf("debug read: %w", err)
	}
	defer f.Close()

	meta, err := ReadOnlyMeta(f)
	if err != nil {
		return nil, fmt.Errorf("read cache %q: %w", path, err)
	}

	switch keeper(meta) {
	case DebugKeepIgnore:
		return nil, nil
	case DebugKeepIndex:
		return meta, nil
	}

	begin := make([]byte, 4096)
	n, _ := f.Read(begin[:])
	meta.ResponseBody = begin[:n]

	return meta, nil
}

func debugPrint(meta *Meta) []byte {
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

	buff.Write(meta.ResponseBody)

	return buff.Bytes()
}

//go:embed debug.css
var debugcss render.H

//go:embed debug.js
var debugjs render.H

func debugRender(index []*Meta) []byte {
	return render.Merge(render.N("html",
		render.N("head",
			render.H(`<meta charset=utf-8>`),
			render.H(`<meta name=viewport content="width=device-width,initial-scale=1">`),
			render.N("title", "Index of HTTP cache"),
			render.N("style", debugcss),
		),
		render.N("body",
			render.H(`<input id=s type=search placeholder="Search ...">`),
			render.N("ul", render.S(index, "", func(meta *Meta) render.Node {
				return render.N("li",
					render.Na("code", "data-id", meta.Path()).N(meta.ID()[:7]),
					" ",
					render.If(len(meta.ResponseBody) > 0, func() render.Node {
						return render.N("",
							render.Na("a", "data-b64", base64.StdEncoding.EncodeToString(debugPrint(meta))).
								A("data-title", meta.URL.Host+"/"+meta.ID()[:7]).
								N("~>"),
							" ",
						)
					}),
					render.IfS(len(meta.ResponseBody) == 0, "-- "),
					render.IfS(meta.Status/100 == 2, render.N("span.s2", "s", meta.Status)),
					render.IfS(meta.Status/100 == 3, render.N("span.s3", "s", meta.Status)),
					render.IfS(meta.Status/100 == 4, render.N("span.s4", "s", meta.Status)),
					render.IfS(meta.Status/100 == 5, render.N("span.s5", "s", meta.Status)),
					" ",
					render.N("span.m", meta.Method),
					" ",
					meta.URL,
					render.N(""),
				)
			})),
			render.N("script", debugjs),
		),
	))
}
