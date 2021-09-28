// Package lib generates starlark documentation from go code
package lib

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Docs is a sortable slice of Doc pointers
type Docs []*Doc

// Len implements the sort.Sortable interface
func (d Docs) Len() int { return len(d) }

// Less implements the sort.Sortable interface
func (d Docs) Less(i, j int) bool { return d[i].Path+d[i].Name < d[j].Path+d[j].Name }

// Swap implements the sort.Sortable interface
func (d Docs) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// Sort sorts all sortable fields in all docs, and the docs list itself
func (d Docs) Sort() {
	for _, doc := range d {
		doc.Sort()
	}
	sort.Sort(d)
}

// Doc is is a documentation document
type Doc struct {
	cfg         config
	Name        string
	Path        string
	Description string
	Functions   Functions
	Types       Types
}

// Sort sorts all sortable fields in the document
func (d *Doc) Sort() {
	if d.cfg.alphaSortFuncs {
		sort.Sort(d.Functions)
	}
	if d.cfg.alphaSortTypes {
		sort.Sort(d.Types)
	}
}

// MarshalIndent writes doc to a string with depth & prefix
func (d *Doc) MarshalIndent(depth int, prefix string) ([]byte, error) {
	buf := &bytes.Buffer{}
	if d.Name == "" {
		buf.WriteString(strings.Repeat(prefix, depth) + DocumentTok.String() + "\n")
	} else {
		buf.WriteString(fmt.Sprintf("%s%s: %s\n", strings.Repeat(prefix, depth), DocumentTok.String(), d.Name))
	}
	if d.Path != "" {
		depth++
		buf.WriteString(strings.Repeat(prefix, depth) + PathTok.String() + ": " + d.Path + "\n")
		depth--
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

// Functions is a sortable slice of Function pointers
type Functions []*Function

// Len implements the sort.Sortable interface
func (f Functions) Len() int { return len(f) }

// Less implements the sort.Sortable interface
func (f Functions) Less(i, j int) bool { return f[i].Signature < f[j].Signature }

// Swap implements the sort.Sortable interface
func (f Functions) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

// Function documents a starlark function
type Function struct {
	FuncName    string
	Receiver    string // should be set by parsing context
	Signature   string
	Description string
	Params      []*Param
	Return      string
}

// Param is an argument to a function
type Param struct {
	Name        string
	Type        string
	Optional    bool
	Description string
}

// Types is a sortable slice of Type pointers
type Types []*Type

// Len implements the sort.Sortable interface
func (t Types) Len() int { return len(t) }

// Less implements the sort.Sortable interface
func (t Types) Less(i, j int) bool { return t[i].Name < t[j].Name }

// Swap implements the sort.Sortable interface
func (t Types) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

// Type documents a constructed type
type Type struct {
	Name        string
	Description string
	Methods     Functions
	Fields      []*Field
	Operators   []*Operator
}

// Sort sorts a Type pointer's Methods
func (t *Type) Sort() {
	sort.Sort(t.Methods)
}

// Field is a property of a constructed Type
type Field struct {
	Name        string
	Type        string
	Description string
}

// Operator documents boolean operation on a constructed type
type Operator struct {
	Opr         string
	Description string
}
