package sch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailAdress(t *testing.T) {
	assert.NoError(t, AnyMail().Match("conduct@golang.org"))
	assert.Error(t, AnyMail().Match(1))
	assert.Error(t, AnyMail().Match(""))
	assert.Error(t, AnyMail().Match("<conduct@golang.org>"))
	assert.Error(t, AnyMail().Match("GO <conduct@golang.org>"))
	assert.Equal(t, `<span class=sch-base title="Email address: user@host">mail-address</span>`, genHTML(AnyMail()))
	assert.Equal(t, "mail-address", AnyMail().String())
}

func TestMime(t *testing.T) {
	assert.Panics(t, func() { MIME("; charset=utf-8") })
	assert.Panics(t, func() { MIME("x; charset=utf-8") })

	assert.NoError(t, MIME("text/html;charset=utf-8").Match("text/html;charset=utf-8"))
	assert.NoError(t, MIME("text/html").Match("text/html"))
	assert.NoError(t, MIME("text/*").Match("text/plain"))
	assert.NoError(t, MIME("*/*").Match("image/png"))

	assert.Error(t, MIME("text/*").Match(1))
	assert.Error(t, MIME("text/*").Match(";charset=utf-8"))
	assert.Error(t, MIME("text/*").Match("text/html;charset=utf-8"))
	assert.Error(t, MIME("text/html;charset=utf-8").Match("text/html;charset=us-ascii"))

	assert.Equal(t, `<span class=sch-xstr>mime(<u>text/html; charset=utf-8</u>)</span>`,
		genHTML(MIME("text/html;charset=utf-8")))
}

func TestRegexp(t *testing.T) {
	assert.NoError(t, Regexp(`^ECI\(\d{4}\)\d{6}$`).Match("ECI(2024)000008"))
	assert.Error(t, Regexp(`^ECI\(\d{4}\)\d{6}$`).Match("ECI(2024)0000008"))
	assert.Error(t, Regexp(`^ECI\(\d{4}\)\d{6}$`).Match(1))
	assert.Equal(t, `<span class=sch-xstr>regexp(<u>^ECI\(\d{4}\)\d{6}$</u>)</span>`, genHTML(Regexp(`^ECI\(\d{4}\)\d{6}$`)))
}

func TestTime(t *testing.T) {
	assert.NoError(t, Time("2006-01-02 15:04:05").Match("2024-10-31 00:45:39"))
	assert.Error(t, Time("2006-01-02 15:04:05").Match("2024-10-31T00:45:39"))
	assert.Error(t, Time("2006-01-02 15:04:05").Match(1))

	assert.Equal(t, `<span class=sch-xstr title="A time value encoded into a string">time(<u>2006-01-02 15:04:05</u>)</span>`,
		genHTML(Time("2006-01-02 15:04:05")))
}

func TestAnyURL(t *testing.T) {
	assert.NoError(t, AnyURL().Match("https://sniffle.eu/dir/file.txt?a=1"))
	assert.Error(t, AnyURL().Match(1))
	assert.Equal(t, `<span class=sch-base title="A HTTP(S) url into a string">url</span>`, genHTML(AnyURL()))
}

func TestURL(t *testing.T) {
	assert.Panics(t, func() { URL("://sniffle.eu/") })

	u := URL("http.s://**.europa.eu/**?a=1&x=*")
	assert.NoError(t, u.Match("https://ec.europa.eu/yolo?a=1&x=42"))
	assert.NoError(t, URL("https://**europa.eu/").Match("https://ec.europa.eu/"))
	assert.NoError(t, URL("https://**europa.eu/").Match("https://europa.eu/"))
	assert.Error(t, u.Match(1))
	assert.Error(t, u.Match("://europa.eu/"))
	assert.Error(t, u.Match("https://user@europa.eu/"))
	assert.Error(t, u.Match("malto:user@europa.eu"))
	assert.Error(t, URL("https://europa.eu/").Match("http://europa.eu/"))
	assert.Error(t, u.Match("https://europa.eu/yolo?a=1&x=42"))
	assert.Error(t, u.Match("https://ec.europa.eu/yolo?a=1&a=2&x=42"))
	assert.Error(t, u.Match("https://ec.europa.eu/yolo?a=2&x=42"))
	assert.Error(t, u.Match("https://ec.europa.eu/yolo?a=1&x=42&b=2"))

	assert.Equal(t, `<span class=sch-xstr>url(<u>http.s://**.europa.eu/**?a=1&amp;x=*</u>)</span>`, genHTML(u))
	assert.Equal(t, "http.s://**.europa.eu/**?a=1&x=*", u.String())
}
