package sch

import "sniffle/tool/render"

type Type interface {
	// Check if a value math the pattern.
	Match(v any) error
	// Print node to generate HTML help page.
	HTML(indent string) render.Node
}

type TypeStringer interface {
	Type
	// Return a string, used in structure field key.
	String() string
}
