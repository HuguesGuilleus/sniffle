package render

import (
	"cmp"
	"fmt"
	"html"
	"html/template"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Used to indicate that this string is already escaped, or the HTML content is safe.
type H = template.HTML

// A int display without space for thousand
type Int int64

// A fake time zone to indicate that this thime, is a date.
// So hour, minute, seconds and milisecond must be ignored.
var DateZone = time.FixedZone("DATE", 0)
var ShortDateZone = time.FixedZone("SHORTDATE", 0)

// Attributes, pair of key value.
// If value is empty, create a empty attribute.
// The value is auto escaped.
type attributes = [][2]string

type NodeAttrBuilder struct {
	tags string
	attr attributes
}

// Create a node builde to add some attributes then create a Node
//
// Usage: Na("a.class1.class2", "href", "https://github.com/").A("title", "GitHub website").N(N("span", "GitHub ..."))
func Na(tags string, key, value string) NodeAttrBuilder {
	return NodeAttrBuilder{tags, attributes{{key, value}}}
}
func (na NodeAttrBuilder) A(key, value string) NodeAttrBuilder {
	na.attr = append(na.attr, [2]string{key, value})
	return na
}
func (na NodeAttrBuilder) N(children ...any) Node {
	return Node{na.tags, na.attr, children}
}

// One HTML node, used to
type Node struct {
	tags string
	attr attributes
	// Children element, maybe alredy escaped string,
	// or something that will be auto escaped.
	children []any
}

// Create a Node.
// Tags pattern is: tagName.class1.class2.class3
// If tags == "", to print only the children without tag.
func N(tags string, children ...any) Node { return Node{tags, nil, children} }

// A zero node who output nothing.
var Z = Node{"", nil, nil}

// If b is true, return f call else return Z.
func If(b bool, f func() Node) Node {
	if b {
		return f()
	}
	return Z
}

// Return b if b is true, else return Z.
func IfS(b bool, n Node) Node {
	if b {
		return n
	}
	return Z
}

// If b is true, return yes call else return no call.
func IfElse(b bool, yes, no func() Node) Node {
	if b {
		return yes()
	} else {
		return no()
	}
}

// Merge all nodes to a HTML page.
// So the root should be a HTML tag.
func Merge(root Node) []byte {
	h := make([]byte, 0)
	h = append(h, `<!DOCTYPE html>`...)
	return root.mergeSlice(h)
}
func (node *Node) mergeSlice(h []byte) []byte {
	if node.tags == "" {
		for _, child := range node.children {
			h = renderChild(h, child)
		}
		return h
	}

	tagsSplited := strings.Split(node.tags, ".")
	tagName := tagsSplited[0]

	// Opening tag
	h = append(h, '<')
	h = append(h, tagName...)
	tagsSplited = tagsSplited[1:]
	if len(tagsSplited) == 1 {
		h = append(h, ` class=`...)
		h = append(h, tagsSplited[0]...)
	} else if len(tagsSplited) > 1 {
		h = append(h, ` class="`...)
		for _, class := range tagsSplited {
			h = append(h, class...)
			h = append(h, ' ')
		}
		h[len(h)-1] = '"'
	}
	for _, attr := range node.attr {
		h = append(h, ' ')
		h = append(h, attr[0]...)
		if v := attr[1]; v == "" {
			// Nothing
		} else if !strings.ContainsAny(v, "<>=\"'`& \t") {
			h = append(h, '=')
			h = append(h, v...)
		} else {
			h = append(h, '=', '"')
			for _, c := range []byte(v) {
				if c == '"' {
					h = append(h, `&#34;`...)
				} else if c == '&' {
					h = append(h, `&amp;`...)
				} else {
					h = append(h, c)
				}
			}
			h = append(h, '"')
		}
	}
	if h[len(h)-1] == '/' {
		h = append(h, ' ')
	}
	h = append(h, '>')

	// Children
	for _, child := range node.children {
		h = renderChild(h, child)
	}

	// End tag
	switch tagName {
	// Source: https://html.spec.whatwg.org/multipage/syntax.html#elements-2
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "source", "track", "wbr":
		// Do not close
	default:
		h = append(h, '<', '/')
		h = append(h, tagName...)
		h = append(h, '>')
	}

	return h
}
func renderChild(h []byte, child any) []byte {
	switch child := child.(type) {
	case NodeAttrBuilder:
		n := child.N()
		h = n.mergeSlice(h)
	case Node:
		h = child.mergeSlice(h)
	case []Node:
		for _, subChild := range child {
			h = subChild.mergeSlice(h)
		}
	case H:
		h = append(h, child...)
	case []H:
		for _, subChild := range child {
			h = append(h, subChild...)
		}
	case string:
		h = append(h, html.EscapeString(child)...)
	case Int:
		h = strconv.AppendInt(h, int64(child), 10)
	case int:
		if child < 0 {
			h = append(h, '-')
			child = -child
		}
		h = renderUint64(h, uint64(child))
	case uint:
		h = renderUint64(h, uint64(child))
	case time.Time:
		h = append(h, `<time datetime=`...)
		if l := child.Location(); l == DateZone {
			h = child.AppendFormat(h, `2006-01-02>2006-01-02`)
		} else if l == ShortDateZone {
			h = child.AppendFormat(h, `2006-01-02>2006_01_02`)
		} else {
			child = child.UTC().Truncate(time.Second)
			h = child.AppendFormat(h, `2006-01-02T15:04:05Z`)
			h = append(h, `>`...)
			h = child.AppendFormat(h, `2006-01-02 15:04:05 UTC`)
		}
		h = append(h, `</time>`...)
	case nil:
		// Nothing
	case []any:
		for _, subChild := range child {
			h = renderChild(h, subChild)
		}
	default:
		h = append(h, fmt.Sprint(child)...)
	}
	return h
}
func renderUint64(h []byte, u uint64) []byte {
	if u >= 1000 {
		h = renderUint64(h, u/1000)
		u %= 1000
		h = append(h, 0xE2, 0x80, 0xAF,
			'0'+byte(u/100),
			'0'+byte(u/10%10),
			'0'+byte(u%10),
		)
	} else {
		if u >= 100 {
			h = append(h, '0'+byte(u/100))
		}
		if u >= 10 {
			h = append(h, '0'+byte(u/10%10))
		}
		h = append(h, '0'+byte(u%10))
	}
	return h
}

func Map[K cmp.Ordered, V any](m map[K]V, f func(k K, v V) Node) []Node {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return mapCall(keys, m, f)
}

func MapReverse[K cmp.Ordered, V any](m map[K]V, f func(k K, v V) Node) []Node {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	slices.Reverse(keys)
	return mapCall(keys, m, f)
}

func mapCall[K cmp.Ordered, V any](keys []K, m map[K]V, f func(k K, v V) Node) []Node {
	children := make([]Node, len(keys))
	for i, k := range keys {
		children[i] = f(k, m[k])
	}
	return children
}

func S[V any](s []V, separator H, f func(v V) Node) []Node {
	if len(s) == 0 {
		return nil
	}
	sep := Z
	if separator != "" {
		sep = N("", separator)
	}
	nodes := make([]Node, 0, len(s)*2-1)
	nodes = append(nodes, f(s[0]))
	for _, v := range s[1:] {
		nodes = append(nodes, sep, f(v))
	}
	return nodes
}

func S2[V any](s []V, separator H, f func(i int, v V) Node) []Node {
	if len(s) == 0 {
		return nil
	}
	sep := Z
	if separator != "" {
		sep = N("", separator)
	}
	nodes := make([]Node, 0, len(s)*2-1)
	nodes = append(nodes, f(0, s[0]))
	for i, v := range s[1:] {
		nodes = append(nodes, sep, f(i+1, v))
	}
	return nodes
}
