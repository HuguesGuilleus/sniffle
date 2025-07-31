// release service for /release/ static pages.
package release

import (
	"time"

	"github.com/HuguesGuilleus/sniffle/common/language"
	"github.com/HuguesGuilleus/sniffle/front/component"
	"github.com/HuguesGuilleus/sniffle/tool"
	"github.com/HuguesGuilleus/sniffle/tool/render"
)

var steps = []render.Node{
	step("2024-07-27", "", "tool", "Begin tool development for this website."),
	step("2024-08-11", "2025-03-26", "eu/ec/eci", "Creation of European Citizens' Initiative pages."),
	step("2024-12-27", "", "release", "Creation of release pages."),
}

func Do(t *tool.Tool) {
	l := language.AllEnglish
	t.WriteFile("/release/index.html", render.Merge(render.Na("html", "lang", l.String()).N(
		render.N("head",
			component.HeadBegin,
			render.N("title", "Release"),
			render.Na("meta", "name", "description").A("content", "Release of this website"),
		),
		render.N("body.edito",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerID", component.HomeAnchor(l), render.Na("a", "href", "..").N("release")),
				),
				render.N("div.headerTitle", "List of change in this website"),
			),
			render.N("div.w",
				render.N("div.timeLine", steps),
			),
			component.Footer(l, 0),
		),
	)))
}

func step(begin, end string, tag string, children ...any) render.Node {
	return render.N("div.timePoint",
		render.N("div.timeHead",
			render.N("span.tag", tag),
			date(begin), render.If(end != "", func() render.Node { return render.N("", " ~> ", date(end)) }),
		),
		render.N("", children...),
	)
}

// Create a date from time format "YYYY-MM-DD"
func date(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.ParseInLocation(time.DateOnly, s, render.DateZone)
	if err != nil {
		panic(err)
	}
	return t
}
