package sch

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"sniffle/tool/render"
	"strings"
	"time"
)

/* EMAIL */

type anyMailType struct{}

func AnyMail() TypeStringer { return anyMailType{} }

func (anyMailType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	if strings.Contains(s, "<") {
		return errors.New("support only simple mail address")
	}

	_, err := mail.ParseAddress(s)
	if err != nil {
		return err
	}

	return nil
}

func (anyMailType) HTML(_ string) render.Node {
	return render.Na(baseMarkup, "title", "Email address: user@host").N("mail-address")
}

func (anyMailType) String() string { return "mail-address" }

/* REGEXP */

type regexpType struct {
	*regexp.Regexp
}

func Regexp(str string) TypeStringer { return regexpType{regexp.MustCompile(str)} }

func (r regexpType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	if !r.Regexp.MatchString(s) {
		return fmt.Errorf("rexpexp /%s/ not match %q", r.Regexp, s)
	}

	return nil
}

func (r regexpType) HTML(_ string) render.Node {
	return render.N(xstringMarkup, "regexp(", render.N("u", r.String()), ")")
}

/* TIME */

type timeType struct {
	layout string
}

// Check that the value is a string formated with layout.
func Time(layout string) Type { return timeType{layout} }

func (t timeType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	_, err := time.Parse(t.layout, s)
	return err
}

func (t timeType) HTML(_ string) render.Node {
	return render.Na(xstringMarkup, "title", "A time value encoded into a string").N(
		"time(", render.N("u", t.layout), ")",
	)
}

/* ANY URL */

type anyUrlType struct{}

func AnyURL() Type { return anyUrlType{} }

func (anyUrlType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	_, err := url.Parse(s)
	return err
}

func (anyUrlType) HTML(_ string) render.Node {
	return render.Na(baseMarkup, "title", "A HTTP(S) url into a string").N("url")
}

/* URL */

type urlType struct {
	url.URL
	eqHost func(string, string) bool
	eqPath func(string, string) bool
	query  url.Values
	rawURL string
}

// URL: sheme://**host/path**?a=1
// scheme can be "http.s" (for http or https), "http" or "https".
//
// For host, "*host" => "www.host" | "dl.host" ... | "host".
// For path, "p**" => any path with prefix is "p".
// For host and path, ** are optional.
//
// "?" => no query
//
// It do no accept URL user info.
//
// If the url syntax is wrong, it panic.
func URL(rawURL string) TypeStringer {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	eqHost := strings.EqualFold
	if strings.HasPrefix(u.Host, "**.") {
		u.Host = strings.TrimPrefix(u.Host, "**")
		eqHost = strings.HasSuffix
	} else if strings.HasPrefix(u.Host, "**") {
		u.Host = strings.TrimPrefix(u.Host, "**")
		eqHost = func(s, t string) bool {
			return s == t || (strings.HasSuffix(s, t) && s[len(s)-len(t)-1] == '.')
		}
	}
	eqPath := strings.EqualFold
	if strings.HasSuffix(u.Path, "**") {
		u.Path = strings.TrimSuffix(u.Path, "**")
		eqPath = strings.HasPrefix
	}

	query := (url.Values)(nil)
	if u.RawQuery != "" || u.ForceQuery {
		query = u.Query()
	}

	return urlType{
		URL:    *u,
		eqHost: eqHost,
		eqPath: eqPath,
		query:  query,
		rawURL: rawURL,
	}
}

func (ut urlType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	u, err := url.Parse(s)
	if err != nil {
		return err
	} else if u.User != nil {
		return errors.New("userinfo in url")
	}

	if ut.Scheme == "http.s" {
		if u.Scheme != "http" && u.Scheme != "https" {
			return fmt.Errorf("wrong protocol %q", s)
		}
	} else if ut.Scheme != u.Scheme {
		return fmt.Errorf("wrong protocol %q, expected %s", s, ut.Scheme)
	}

	if !ut.eqHost(u.Host, ut.Host) || !ut.eqPath(u.Path, ut.Path) {
		return fmt.Errorf("wrong URL %q expected %q", s, ut.rawURL)
	}

	if ut.query != nil {
		q := u.Query()
		for k, v := range ut.query {
			if notSameQuery(v, q[k]) {
				return fmt.Errorf("url.query(%q) != %q", ut.query, u.Query())
			}
			delete(q, k)
		}
		if len(q) != 0 {
			return fmt.Errorf("url.query(%q) != %q", ut.query, u.Query())
		}
	}

	return nil
}
func notSameQuery(expected, recieved []string) bool {
	if expected[0] == "*" {
		return false
	}
	if len(expected) != len(recieved) {
		return true
	}
	for i, e := range expected {
		if e != recieved[i] {
			return true
		}
	}
	return false
}

func (ut urlType) HTML(_ string) render.Node {
	return render.N(xstringMarkup, "url(", render.N("u", ut.rawURL), ")")
}

func (ut urlType) String() string { return ut.rawURL }
