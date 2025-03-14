package release

import (
	"fmt"
	"sniffle/common/language"
	"sniffle/front/component"
	"sniffle/tool"
	"sniffle/tool/render"
	"strconv"
	"time"
)

var steps = map[render.Int][]render.Node{
	2024: {
		step("2024-08-11", "", "eu/ec/eci", "Creation of European Citizens' Initiative pages and development tools to built this website."),
		step("2024-12-27", "", "release", "Creation of release pages."),
	},
}

func Do(t *tool.Tool) {
	basePath := "/release/"
	hostURL := t.HostURL + basePath
	l := language.AllEnglish

	t.WriteFile(basePath+"index.html", render.Merge(render.Na("html", "lang", l.String()).N(
		component.Head(l, hostURL, "Release index", "Index of all release"),
		render.N("body.edito",
			component.TopHeader(l),
			render.N("header",
				render.N("div.headerSup",
					render.N("div.headerId", component.HomeAnchor(l), render.Na("a", "href", ".").N("release")),
				),
				render.N("div.headerTitle", "List of change in this website"),
			),
			render.N("ul.w.home",
				render.Map(steps, func(year render.Int, _ []render.Node) render.Node {
					return render.N("li", render.Na("a", "href", fmt.Sprintf("%d/", year)).N(year))
				}),
			),
			component.Footer(l, 0),
		),
	)))

	for year, steps := range steps {
		y := strconv.Itoa(int(year))
		t.WriteFile(fmt.Sprintf("%s%d/index.html", basePath, year), render.Merge(render.Na("html", "lang", l.String()).N(
			component.Head(l, hostURL, y+" release", "Release of "+y),
			render.N("body.edito",
				component.TopHeader(l),
				render.N("header",
					render.N("div.headerSup",
						render.N("div.headerId", component.HomeAnchor(l), render.Na("a", "href", "..").N("release")),
						render.N("div.headerId", year),
					),
					render.N("div.headerTitle", "List of change in this website in ", year),
				),
				render.N("div.w",
					render.N("div.timeLine", steps),
				),
				component.Footer(l, 0),
			),
		)))
	}
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
