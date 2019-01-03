package lib

import (
	"fmt"
	"io"
)

// Parse consumes a reader of outline data, creating a outline document
func Parse(r io.Reader) (doc *Doc, err error) {
	p := parser{s: NewScanner(r)}
	doc, err = p.read()
	if err == io.EOF {
		err = nil
	}

	return
}

// parser is a state machine for serializing a documentation struct from a byte stream
type parser struct {
	s *Scanner

	buf struct {
		tok          Token
		line, indent int
		n            int
	}

	line   int
	indent int // indentation level of current line
}

func (p *parser) scan() (tok Token) {
	if p.buf.n > 0 {
		// fmt.Printf("using scan buffer %s %#v\n", p.buf.tok.Type, p.buf.tok)
		tok = p.buf.tok
		p.indent = p.buf.indent
		p.line = p.buf.line
		p.buf.n = 0
		return
	}

	defer func() {
		p.buf.tok = tok
		p.buf.line = p.line
		p.buf.indent = p.indent
	}()

	for {
		tok = p.s.Scan()
		switch tok.Type {
		case NewlineTok:
			p.indent = 0
			p.line++
		case IndentTok:
			p.indent++
		case EofTok:
			return
		default:
			// fmt.Printf("returning token: %s %#v\n", tok.Type, tok.Text)
			return
		}
	}
}

func (p *parser) unscan() {
	p.buf.n = 1
}

func (p *parser) read() (doc *Doc, err error) {
	for {
		tok := p.scan()
		switch tok.Type {
		case DocumentTok:
			doc, err = p.readDocument(p.indent)
			return
		case EofTok:
			return
		}
	}
}

func (p *parser) readDocument(baseIndent int) (doc *Doc, err error) {
	doc = &Doc{}
	tok := p.scan()
	if tok.Type == TextTok {
		doc.Name = tok.Text
	} else {
		p.unscan()
	}

	for {
		tok := p.scan()
		switch tok.Type {
		case FunctionsTok:
			if doc.Functions, err = p.readFunctions(baseIndent + 1); err != nil {
				return
			}
		case TypesTok:
			if doc.Types, err = p.readTypes(baseIndent + 1); err != nil {
				return
			}
		default:
			p.unscan()
			return
		}
	}
}

func (p *parser) readFunctions(baseIndent int) (funcs []*Function, err error) {
	for {
		var fn *Function
		if fn, err = p.readFunction(baseIndent + 1); err != nil || fn == nil {
			return
		}
		funcs = append(funcs, fn)
	}
}

func (p *parser) readFunction(baseIndent int) (fn *Function, err error) {
	// read signature
	tok := p.scan()
	if p.indent < baseIndent || tok.Type != TextTok {
		p.unscan()
		return
	}

	fn = &Function{Signature: tok.Text}
	fn.Description, err = p.readText(baseIndent + 1)
	return
}

func (p *parser) readTypes(baseIndent int) (types []*Type, err error) {
	for {
		var t *Type
		if t, err = p.readType(baseIndent + 1); err != nil || t == nil {
			return
		}
		types = append(types, t)
	}
}

func (p *parser) readType(baseIndent int) (t *Type, err error) {
	// read signature
	tok := p.scan()
	if p.indent < baseIndent || tok.Type != TextTok {
		p.unscan()
		return
	}

	t = &Type{Name: tok.Text}

	for {
		tok = p.scan()
		if p.indent <= baseIndent {
			p.unscan()
			return
		}

		switch tok.Type {
		case FieldsTok:
			t.Fields, err = p.readFields(baseIndent + 1)
			if err != nil {
				return
			}
		case OperatorsTok:
			t.Operators, err = p.readOperators(baseIndent + 1)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("unexpexted token: %s: %s %d %d", tok.Type, tok.Text, p.indent, baseIndent)
			return

			// TODO (b5): handle field descriptions
			// case TextTok:
			// fmt.Println("text?")
			// p.unscan()
			// t.Description, err = p.readText(baseIndent + 1)
		}
	}
}

func (p *parser) readFields(baseIndent int) (fields []*Field, err error) {
	for {
		var f *Field
		if f, err = p.readField(baseIndent + 1); err != nil || f == nil {
			return
		}
		fields = append(fields, f)
	}
}

func (p *parser) readField(baseIndent int) (field *Field, err error) {
	tok := p.scan()
	if p.indent < baseIndent || tok.Type != TextTok {
		p.unscan()
		return
	}
	field = &Field{Name: tok.Text}
	return
}

func (p *parser) readOperators(baseIndent int) (ops []*Operator, err error) {
	for {
		var o *Operator
		if o, err = p.readOperator(baseIndent + 1); err != nil || o == nil {
			return
		}
		ops = append(ops, o)
	}
}

func (p *parser) readOperator(baseIndent int) (op *Operator, err error) {
	tok := p.scan()
	if p.indent < baseIndent || tok.Type != TextTok {
		p.unscan()
		return
	}
	op = &Operator{Opr: tok.Text}
	return
}

func (p *parser) readText(baseIndent int) (str string, err error) {
	for {
		tok := p.scan()
		// fmt.Printf("readText %d %d %s %#v\n", baseIndent, p.indent, tok.Type, tok)
		if p.indent < baseIndent || tok.Type != TextTok {
			p.unscan()
			return
		}

		if str == "" {
			str = tok.Text
		} else {
			str += " " + tok.Text
		}
	}
}

func (p *parser) errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
