package securehtml

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/HuguesGuilleus/sniffle/tool/render"

	"golang.org/x/net/html"
)

func ParseURL(s string) *url.URL {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	u, _ := url.Parse(s)
	if u == nil {
		return nil
	}
	switch u.Scheme {
	case "https", "http":
		return u
	case "":
		return ParseURL("https://" + s)
	default:
		return nil
	}
}

func Text(src string, limit int) string {
	root, err := html.Parse(strings.NewReader(src))
	if err != nil {
		return src
	}
	buff := bytes.Buffer{}
	plainText(root, &buff)

	data := buff.Bytes()
	whitespace := false
	i := 0
	for _, b := range data {
		switch b {
		case '\f', '\r', '\n', '\t', ' ':
			whitespace = true
		default:
			if whitespace {
				whitespace = false
				data[i] = ' '
				i++
			}
			data[i] = b
			i++
		}
	}

	begin := 0
	if len(data) > 0 && data[0] == ' ' {
		begin = 1
	}

	s := string(data[begin:i])
	nb := 0
	for i := range s {
		nb++
		if nb > limit {
			return s[:i]
		}
	}

	return s
}
func plainText(node *html.Node, buff *bytes.Buffer) {
	if node.Type == html.TextNode {
		buff.WriteString(node.Data)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		plainText(child, buff)
	}
}

func TextWithURL(s string) render.H {
	buff := buffer{}
	buff.writeStringWithUrl(s)
	return render.H(buff.buffer.String())
}
