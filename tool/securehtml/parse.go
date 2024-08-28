package securehtml

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"sniffle/tool/render"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var urlRegexp = regexp.MustCompile(`(https?\://\S+)`)

type buffer struct {
	buffer bytes.Buffer
	// Deep pre, code
	keepSpace int

	spaceForbiden bool
}

func isBlock(tag atom.Atom) bool {
	switch tag {
	case atom.P, atom.Br, atom.Blockquote, atom.Ul, atom.Ol, atom.Li, atom.Pre:
		return true
	default:
		return false
	}
}

func keepSpaceTag(tag atom.Atom) bool {
	switch tag {
	case atom.Pre, atom.Code:
		return true
	default:
		return false
	}
}

func isWiteSpace(b byte) bool {
	switch b {
	case '\x00', ' ', '\t', '\f', '\r', '\n':
		return true
	default:
		return false
	}
}

// Remove the whitespace at end of the buffer.
func (buff *buffer) trimWhiteSpaceEnd() {
	for buff.buffer.Len() > 0 {
		if l := buff.buffer.Len(); isWiteSpace(buff.buffer.Bytes()[l-1]) {
			buff.buffer.Truncate(l - 1)
		} else {
			return
		}
	}
}

func (buff *buffer) WriteString(s string) {
	if buff.keepSpace > 0 {
		buff.spaceForbiden = false
		buff.buffer.WriteString(html.EscapeString(s))
		return
	}

	space := buff.spaceForbiden
	for _, b := range []byte(s) {
		if isWiteSpace(b) {
			if !space {
				buff.buffer.WriteByte(' ')
			}
			space = true
		} else {
			space = false
			switch b {
			case '<':
				buff.buffer.WriteString(`&gt;`)
			case '>':
				buff.buffer.WriteString(`&lt;`)
			case '&':
				buff.buffer.WriteString(`&amp;`)
			case '"':
				buff.buffer.WriteString(`&#34;`)
			case '\'':
				buff.buffer.WriteString(`&#39;`)
			default:
				buff.buffer.WriteByte(b)
			}
		}
	}

	buff.spaceForbiden = space
}

func (buff *buffer) writeStringWithUrl(s string) {
	loc := urlRegexp.FindStringIndex(s)
	if loc == nil {
		buff.WriteString(s)
		return
	}

	buff.WriteString(s[:loc[0]])

	sub := s[loc[0]:loc[1]]
	if u, _ := url.Parse(sub); u == nil {
		buff.WriteString(sub)
	} else {
		buff.addAnchor(u)
		buff.WriteString(sub)
		buff.end(atom.A)
	}

	buff.writeStringWithUrl(s[loc[1]:])
}

func (buff *buffer) addAnchor(href *url.URL) {
	buff.buffer.WriteString(`<a href="`)
	buff.buffer.WriteString(html.EscapeString(href.String()))
	buff.buffer.WriteString(`">`)
}

func (buff *buffer) add(tag atom.Atom) {
	if keepSpaceTag(tag) {
		buff.keepSpace++
	}
	if isBlock(tag) {
		buff.trimWhiteSpaceEnd()
		buff.spaceForbiden = true
	}
	buff.buffer.WriteByte('<')
	buff.buffer.WriteString(tag.String())
	buff.buffer.WriteByte('>')
}

func (buff *buffer) end(tag atom.Atom) {
	if keepSpaceTag(tag) {
		buff.keepSpace--
	}
	if isBlock(tag) {
		buff.trimWhiteSpaceEnd()
		buff.spaceForbiden = true
	}
	buff.buffer.WriteString("</")
	buff.buffer.WriteString(tag.String())
	buff.buffer.WriteByte('>')
}

func (buff *buffer) walk(node *html.Node) {
	if node.Type == html.TextNode {
		buff.writeStringWithUrl(node.Data)
		return
	} else if node.Type == html.ElementNode {
		if node.Namespace != "" {
			return
		}
		// Inspired from: https://docs.joinmastodon.org/spec/activitypub/#sanitization
		switch node.DataAtom {
		case atom.P, atom.Blockquote, atom.Del, atom.Pre, atom.Code, atom.Em, atom.Strong, atom.B, atom.I, atom.U, atom.Ul, atom.Ol, atom.Li, atom.Sub, atom.Sup:
			buff.add(node.DataAtom)
			defer buff.end(node.DataAtom)
		case atom.A:
			rawURL := ""
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					rawURL = attr.Val
					break
				}
			}
			parsedURL := ParseURL(rawURL)
			if rawURL != "" && parsedURL != nil {
				buff.addAnchor(parsedURL)
				defer buff.end(atom.A)
			} else if rawURL != "" {
				buff.add(atom.Del)
				defer buff.end(atom.Del)
				buff.buffer.WriteString("[")
				buff.WriteString(rawURL)
				buff.buffer.WriteString("] ")
			}

		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
			buff.add(atom.P)
			defer buff.end(atom.P)
			buff.add(atom.Strong)
			defer buff.end(atom.Strong)
		case atom.Html, atom.Body, atom.Span, atom.Div:
			// Ignore markup but render children.
		case atom.Br:
			buff.add(atom.Br)
		case 0, atom.Hr, atom.Iframe, atom.Script, atom.Style, atom.Head:
			return // remove
		default:
			fmt.Printf("securehtml.Secure(), unknown node type: %s (%q)\n", node.DataAtom.String(), node.Data)
			return
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		buff.walk(child)
	}
}

// Secure by escape unknown markup and by removing attribute.
func Secure(src string) render.H {
	root, err := html.Parse(strings.NewReader(src))
	if err != nil {
		return render.H(html.EscapeString(src))
	}

	buff := buffer{}
	buff.walk(root)
	return render.H(buff.buffer.String())
}
