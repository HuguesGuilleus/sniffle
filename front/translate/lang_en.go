package translate

import (
	"html/template"
	"time"
)

type LangEnglish struct{}

func (LangEnglish) DateHourLong(t time.Time) template.HTML {
	return template.HTML(t.Format("<time datetime=2006-01-02T15:04:05Z>Monday, 02 January 2006 15:04:05 MST</time>"))
}

func (LangEnglish) DateShort(t time.Time) template.HTML {
	return template.HTML(t.Format("<time datetime=2006-01-02>2006-01-02</time>"))
}

func (LangEnglish) DateLong(t time.Time) template.HTML {
	return template.HTML(t.Format("<time datetime=2006-01-02>January 2, 2006</time>"))
}
