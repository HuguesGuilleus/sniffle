package translate

import (
	"fmt"
	"html/template"
	"time"
)

type LangFrench struct{}

func (LangFrench) DateHourLong(t time.Time) template.HTML {
	weekday := ""
	switch t.Weekday() {
	case time.Sunday:
		weekday = "dimanche"
	case time.Monday:
		weekday = "lundi"
	case time.Tuesday:
		weekday = "mardi"
	case time.Wednesday:
		weekday = "mercredi"
	case time.Thursday:
		weekday = "jeudi"
	case time.Friday:
		weekday = "vendredi"
	case time.Saturday:
		weekday = "samedi"
	}

	year, month, day := t.Date()
	return template.HTML(fmt.Sprintf("<time datetime=%d-%02d-%02dT%02d:%02d:%02dZ>%s %d %s %d à %02d:%02d:%02d %s</time>",
		year,
		month,
		day,
		t.Hour(),
		t.Minute(),
		t.Second(),

		weekday,
		day,
		LangFrench{}.month(month),
		year,
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Location(),
	))
}

func (LangFrench) DateShort(t time.Time) template.HTML {
	return template.HTML(t.Format("<time datetime=2006-01-02>02/01/2006</time>"))
}

func (LangFrench) DateLong(t time.Time) template.HTML {
	year, month, day := t.Date()
	return template.HTML(fmt.Sprintf("<time datetime=%d-%02d-%02d>%d %s %d</time>",
		year,
		month,
		day,
		day,
		LangFrench{}.month(month),
		year,
	))
}

func (LangFrench) month(month time.Month) string {
	switch month {
	case time.January:
		return "janvier"
	case time.February:
		return "février"
	case time.March:
		return "mars"
	case time.April:
		return "avril"
	case time.May:
		return "mai"
	case time.June:
		return "juin"
	case time.July:
		return "juillet"
	case time.August:
		return "août"
	case time.September:
		return "septembre"
	case time.October:
		return "octobre"
	case time.November:
		return "novembre"
	case time.December:
		return "décembre"
	default:
		return "???"
	}
}
