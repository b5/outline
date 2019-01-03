// Package outline generates starlark documentation from go code
package lib

import (
	"bytes"
	"fmt"
	"strings"
)

// Doc is is a documentation document
type Doc struct {
	Name        string
	Description string
	Functions   []*Function
	Types       []*Type
}

// MarshalIndent writes doc to a string with depth & prefix
func (d *Doc) MarshalIndent(depth int, prefix string) ([]byte, error) {
	buf := &bytes.Buffer{}
	if d.Name == "" {
		buf.WriteString(strings.Repeat(prefix, depth) + DocumentTok.String() + "\n")
	} else {
		buf.WriteString(fmt.Sprintf("%s%s: %s\n", strings.Repeat(prefix, depth), DocumentTok.String(), d.Name))
	}
	if d.Description != "" {
		depth++
		buf.WriteString(strings.Repeat(prefix, depth) + d.Description + "\n")
		depth--
	}
	if d.Functions != nil {
		depth++
		buf.WriteString(strings.Repeat(prefix, depth) + FunctionsTok.String() + ":\n")
		depth++
		for _, fn := range d.Functions {
			buf.WriteString(strings.Repeat(prefix, depth) + fn.Signature + "\n")
			if fn.Description != "" {
				for _, p := range strings.Split(fn.Description, "\n") {
					buf.WriteString(strings.Repeat(prefix, depth+1) + p + "\n")
				}
			}
		}
		depth -= 2
	}

	if d.Types != nil {
		depth++
		buf.WriteString(strings.Repeat(prefix, depth) + TypesTok.String() + ":\n")
		depth++
		for _, t := range d.Types {
			buf.WriteString(strings.Repeat(prefix, depth) + t.Name + "\n")
			if t.Description != "" {
				depth++
				for _, p := range strings.Split(t.Description, "\n") {
					buf.WriteString(strings.Repeat(prefix, depth+1) + p + "\n")
				}
				depth--
			}
			if len(t.Fields) > 0 {
				depth++
				buf.WriteString(strings.Repeat(prefix, depth) + FieldsTok.String() + ":\n")
				depth++
				for _, f := range t.Fields {
					buf.WriteString(strings.Repeat(prefix, depth) + f.Name)
					if f.Type != "" {
						buf.WriteString(" " + f.Type)
					}
					buf.WriteString("\n")
				}
				depth -= 2
			}
		}
		depth -= 2
	}

	return buf.Bytes(), nil
}

// Function documents a starlark function
type Function struct {
	Signature   string
	Description string
}

// Type documents a constructed type
type Type struct {
	Name        string
	Description string
	Methods     []*Function
	Fields      []*Field
	Operators   []*Operator
}

type Field struct {
	Name        string
	Type        string
	Description string
}

// Operator documents
type Operator struct {
	Opr         string
	Description string
}
