package tool

import (
	"html/template"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Secure HTML by escape unknown markup and by removing attribute.
func SecureHTML(src string) template.HTML {
	root, err := html.Parse(strings.NewReader(src))
	if err != nil {
		return template.HTML(html.EscapeString(src))
	}
	buff := strings.Builder{}
	renderSecureHTML(root, &buff)
	return template.HTML(buff.String())
}
func renderSecureHTML(node *html.Node, buff *strings.Builder) {
	if node.Type == html.TextNode {
		buff.WriteString(html.EscapeString(node.Data))
	} else if node.Type == html.ElementNode && node.Namespace != "" {
		return
	} else if node.Type == html.ElementNode {
		// Inspired from: https://docs.joinmastodon.org/spec/activitypub/#sanitization
		switch node.DataAtom {
		case atom.Br, atom.Script, atom.Style:
			return // remove
		case atom.P, atom.Span, atom.Del, atom.Pre, atom.Code, atom.Em, atom.Strong, atom.B, atom.I, atom.U, atom.Ul, atom.Ol, atom.Li, atom.Blockquote:
			buff.WriteByte('<')
			buff.WriteString(node.DataAtom.String())
			buff.WriteByte('>')
			defer func() {
				buff.WriteString("</")
				buff.WriteString(node.DataAtom.String())
				buff.WriteByte('>')
			}()
		case atom.A:
			rawURL := ""
			parsedURL := (*url.URL)(nil)
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					rawURL = attr.Val
					parsedURL = ParseURL(attr.Val)
					break
				}
			}

			if rawURL != "" && parsedURL != nil {
				buff.WriteString(`<a href="`)
				buff.WriteString(html.EscapeString(parsedURL.String()))
				buff.WriteString(`">`)
				defer buff.WriteString("</a>")
			} else if rawURL != "" {
				buff.WriteString(`<del>[`)
				buff.WriteString(html.EscapeString(rawURL))
				buff.WriteString("] ")
				defer buff.WriteString("</del>")
			}

		case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
			buff.WriteString("<p><strong>")
			defer buff.WriteString("</strong></p>")
		default:
			// Ignore markup but render children.
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		renderSecureHTML(child, buff)
	}
}

func ParseURL(s string) *url.URL {
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
	default:
		return nil
	}
}
