package sch

import (
	"errors"
	"fmt"
	"mime"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/HuguesGuilleus/sniffle/tool/render"
)

var (
	justHostRegexp = regexp.MustCompile(`^([\w-]+\.)+[\w-]{2,}$`)
	anyMailRegexp  = regexp.MustCompile(`^\w[\w.-]+\w@([\w-]+\.)+[\w-]{2,}$`)
)

/* EMAIL */

type anyMailType struct{}

func AnyMail() TypeStringer { return anyMailType{} }

func (anyMailType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}

	if !anyMailRegexp.MatchString(s) {
		return fmt.Errorf("not email adress: %q", s)
	}

	return nil
}

func (anyMailType) HTML(_ string) render.Node {
	return render.Na(baseMarkup, "title", "Email address: user@host").N("mail-address")
}

func (anyMailType) String() string { return "mail-address" }

/* Flags */

type flagsType struct {
	flags     map[string]bool
	separator string
	flagsJoin string
}

// A string with multiple flags join by the separator.
func Flags(separator string, flags ...string) TypeStringer {
	flagsSet := make(map[string]bool, len(flags))
	for _, f := range flags {
		flagsSet[f] = true
	}

	return &flagsType{
		flags:     flagsSet,
		separator: separator,
		flagsJoin: strings.Join(flags, separator),
	}
}

func (f *flagsType) String() string { return f.flagsJoin }

func (f *flagsType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}

	alreadyFound := make(map[string]bool)
	for _, sub := range strings.Split(s, f.separator) {
		if !f.flags[sub] {
			return fmt.Errorf("substring %q from %q do not exist in flags %s", sub, s, f.flagsJoin)
		}
		if alreadyFound[sub] {
			return fmt.Errorf("Already found substring %q in %s", sub, s)
		}
		alreadyFound[sub] = true
	}

	return nil
}

func (f *flagsType) HTML(string) render.Node {
	return render.N(xstringMarkup, "flags(",
		render.Map(f.flags, func(k string, _ bool) render.Node {
			return render.N("", render.N(stringMarkup, strconv.Quote(k)), ", ")
		}),
		"separator=", strconv.Quote(f.separator),
		")")
}

/* MIME */

type mimeType struct {
	mainType string
	subType  string
	params   map[string]string
}

func MIME(pattern string) TypeStringer {
	main, sub, params, err := parseMIME(pattern)
	if err != nil {
		panic(err.Error())
	}
	return mimeType{main, sub, params}
}
func parseMIME(s string) (main, sub string, params map[string]string, err error) {
	mediatype := ""
	mediatype, params, err = mime.ParseMediaType(s)
	if err != nil {
		return
	}
	ok := false
	main, sub, ok = strings.Cut(mediatype, "/")
	if !ok {
		return "", "", nil, fmt.Errorf("Wrong media type of %q", s)
	}
	return
}

func (t mimeType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}

	main, sub, params, err := parseMIME(s)
	if err != nil {
		return fmt.Errorf("cannot parse mime type of %q: %w", s, err)
	}
	if (t.mainType != "*" && t.mainType != main) ||
		(t.subType != "*" && t.subType != sub) ||
		notEqualMap(t.params, params) {
		return fmt.Errorf("not same mime %q: get %q", t.String(), s)
	}

	return nil
}
func notEqualMap(a, b map[string]string) bool {
	if len(a) != len(b) {
		return true
	}
	for k, v := range a {
		if v != b[k] {
			return true
		}
	}
	return false
}

func (t mimeType) HTML(_ string) render.Node {
	return render.N(xstringMarkup, "mime(", render.N("u", t.String()), ")")
}

func (t mimeType) String() string {
	return mime.FormatMediaType(t.mainType+"/"+t.subType, t.params)
}

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
		return fmt.Errorf("string %q does not match regexp /%s/", s, r.Regexp)
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

// Any HTTP or HTTPS url without user information.
func AnyURL() Type { return anyUrlType{} }

func (anyUrlType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}
	if anyMailRegexp.MatchString(s) {
		return fmt.Errorf("expected URL, get mail address: %q", s)
	}
	if justHostRegexp.MatchString(s) {
		return nil
	}
	u, err := url.Parse(s)
	if u.Scheme != "http" && u.Scheme != "https" && u.Scheme != "" {
		return fmt.Errorf("not http.s scheme: %q", s)
	}
	if u.User != nil {
		return fmt.Errorf("user information is not accepted: %q", s)
	}
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

	if anyMailRegexp.MatchString(s) {
		return fmt.Errorf("expected URL, get mail address: %q", s)
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

/* UUID */

var uuidRegexp = regexp.MustCompile(`^[a-f\d]{8}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{4}-[a-f\d]{12}$`)

type uuidType struct{}

// A UUID string.
// Match string like`xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` where x is digit or a to f in lower case.
func UUID() Type { return uuidType{} }

func (uuidType) Match(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf(notStringFormat, v)
	}

	if !uuidRegexp.MatchString(s) {
		return fmt.Errorf("string %q is not a uuid", s)
	}

	return nil
}

func (uuidType) HTML(string) render.Node {
	return render.N(xstringMarkup, "uuid")
}
