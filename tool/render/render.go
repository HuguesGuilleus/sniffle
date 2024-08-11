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

// Render safe HTML.
type HTML interface {
	HTML() H
}

// A fake time zone to indicate that this thime, is a date.
// So hour, minute, seconds and milisecond must be ignored.
var DateZone = time.FixedZone("DATE", 0)

// Attributes, pair of key value.
// Key is considerated as same.
// If value is empty, create a empty attribute.
// The value is auto escaped.
type Attributes [][2]string

// Create a new attribute
func A(key, value string) Attributes {
	return Attributes{{key, value}}
}

// Append a new attribute to the slice.
// Can be chained.
func (a Attributes) A(key, value string) Attributes {
	return append(a, [2]string{key, value})
}

// One HTML node, used to
type Node struct {
	tags string
	attr Attributes
	// Children element, maybe alredy escaped string,
	// or something that will be auto escaped.
	children []any
}

// Create a Node.
// Tags pattern is: tagName.class1.class2.class3
// If tags == "!", to print only the children without tag.
func N(tags string, children ...any) Node { return Node{tags, nil, children} }

// Create a Node with options.
func No(tags string, attr Attributes, children ...any) Node {
	return Node{tags, attr, children}
}

// Create a zero node, that production nothing.
var Z = Node{"!", nil, nil}

// If b is true, return f call else return Z.
func If(b bool, f func() Node) Node {
	if b {
		return f()
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
	if node.tags == "!" {
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
	case "link", "img", "meta":
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
	case Node:
		h = child.mergeSlice(h)
	case []Node:
		for _, subChild := range child {
			h = subChild.mergeSlice(h)
		}
	case HTML:
		h = append(h, child.HTML()...)
	case H:
		h = append(h, child...)
	case []H:
		for _, subChild := range child {
			h = append(h, subChild...)
		}
	case string:
		h = append(h, html.EscapeString(child)...)
	case int:
		h = append(h, strconv.Itoa(child)...)
	case uint:
		h = append(h, strconv.FormatUint(uint64(child), 10)...)
	case time.Time:
		h = append(h, `<time datetime=`...)
		if child.Location() == DateZone {
			h = child.AppendFormat(h, `2006-01-02>2006-01-02`)
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

func Map[K cmp.Ordered, V any](m map[K]V, f func(k K, v V) Node) []Node {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	children := make([]Node, len(keys))
	for i, k := range keys {
		children[i] = f(k, m[k])
	}

	return children
}

func Slice[V any](s []V, f func(i int, v V) Node) []Node {
	nodes := make([]Node, len(s))
	for i, v := range s {
		nodes[i] = f(i, v)
	}
	return nodes
}

func SliceSeparator[V any](s []V, separator H, f func(i int, v V) Node) []Node {
	nodes := make([]Node, 0, len(s)*2)
	for i, v := range s {
		if i != 0 {
			nodes = append(nodes, N("!", separator))
		}
		nodes = append(nodes, f(i, v))
	}
	return nodes
}
