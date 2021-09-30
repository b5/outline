package lib

import (
	"fmt"
	"io"
	"strings"
)

// ParseFirst consumes a reader of outline data, creating and returning the first outline document
// it encounters. ParseFirst consumes the entire reader
// TODO(b5): don't consume the entire reader. return after the first complete document
func ParseFirst(r io.Reader, opts ...Option) (doc *Doc, err error) {
	docs, err := Parse(r, opts...)
	if err != nil {
		return nil, err
	}

	if len(docs) > 0 {
		return docs[0], nil
	}
	return nil, nil
}

// Parse consumes a reader of data that contains zero or more outlines
// creating and returning any documents it finds
func Parse(r io.Reader, opts ...Option) (docs Docs, err error) {
	cfg, err := parseOptions(opts)
	if err != nil {
		return docs, err
	}
	p := parser{s: newScanner(r), cfg: cfg}
	for {
		doc, err := p.read()
		if doc == nil && err == nil {
			return docs, nil
		}
		doc.Sort()
		docs = append(docs, doc)
		if err != nil {
			if err == io.EOF {
				docs = append(docs, doc)
				return docs, nil
			}
			return nil, err
		}
	}
}

// parser is a state machine for serializing a documentation struct from a byte stream
type parser struct {
	s   *scanner
	cfg config

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
		case eofTok:
			return
		default:
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
		case eofTok:
			return
		}
	}
}

func (p *parser) readDocument(baseIndent int) (doc *Doc, err error) {
	doc = &Doc{cfg: p.cfg}
	tok := p.scan()
	if tok.Type == TextTok {
		doc.Name = tok.Text
	} else {
		p.unscan()
	}

	for {
		tok := p.scan()
		if p.indent < baseIndent {
			p.unscan()
			return
		}

		switch tok.Type {
		case DocumentTok:
			if p.indent == baseIndent {
				p.unscan()
				return
			}

			err = fmt.Errorf("outline documents cannot be nested")
			return
		case PathTok:
			if doc.Path, err = p.readMultilineText(p.indent); err != nil {
				return
			}
		case FunctionsTok:
			if doc.Functions, err = p.readFunctions(doc.Name, p.indent); err != nil {
				return
			}
		case TypesTok:
			if doc.Types, err = p.readTypes(p.indent); err != nil {
				return
			}
		case TextTok:
			// only read descriptions when indented
			if p.indent > baseIndent {
				p.unscan()
				text, err := p.readMultilineText(p.indent)
				if err != nil {
					return doc, err
				}
				doc.Description = text
			}
		default:
			p.unscan()
			return
		}
	}
}

func (p *parser) readFunctions(receiver string, baseIndent int) (funcs []*Function, err error) {
	for {
		var fn *Function
		if fn, err = p.readFunction(receiver, baseIndent+1); err != nil || fn == nil {
			return
		}
		funcs = append(funcs, fn)
	}
}

func (p *parser) readFunction(receiver string, baseIndent int) (fn *Function, err error) {
	// read signature
	tok := p.scan()
	if p.indent < baseIndent || tok.Type != TextTok {
		p.unscan()
		return
	}

	funcName := ""
	pos := strings.Index(tok.Text, "(")
	if pos != -1 {
		funcName = tok.Text[:pos]
	}

	fn = &Function{FuncName: funcName, Receiver: receiver, Signature: tok.Text}
	for {
		tok := p.scan()
		if p.indent <= baseIndent {
			p.unscan()
			return
		}

		switch tok.Type {
		case ParamsTok:
			if fn.Params, err = p.readParams(p.indent); err != nil {
				return
			}
		case ReturnTok:
			if fn.Return, err = p.readMultilineText(p.indent); err != nil {
				return
			}
		case ExamplesTok:
			if fn.Examples, err = p.readExamples(p.indent); err != nil {
				return fn, err
			}
		case TextTok:
			p.unscan()
			if fn.Description, err = p.readTextBlock(p.indent); err != nil {
				return
			}
		default:
			p.unscan()
			return
		}
	}
}

func (p *parser) readParams(baseIndent int) (params []*Param, err error) {
	for {
		var param *Param
		if param, err = p.readParam(baseIndent + 1); err != nil || param == nil {
			return
		}
		params = append(params, param)
	}
}

func (p *parser) readParam(baseIndent int) (param *Param, err error) {
	tok := p.scan()
	if p.indent < baseIndent || tok.Type != TextTok {
		p.unscan()
		return
	}

	spl := strings.Split(tok.Text, " ")
	switch len(spl) {
	default:
		param = &Param{Name: tok.Text}
	case 2:
		param = &Param{
			Name: spl[0],
			Type: spl[1],
		}
	}

	param.Description, err = p.readMultilineText(baseIndent + 1)
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
			if t.Fields, err = p.readFields(p.indent); err != nil {
				return
			}
		case MethodsTok, FunctionsTok:
			if t.Methods, err = p.readFunctions(t.Name, p.indent); err != nil {
				return
			}
		case OperatorsTok:
			if t.Operators, err = p.readOperators(p.indent); err != nil {
				return
			}
		case TextTok:
			p.unscan()
			if t.Description, err = p.readTextBlock(p.indent); err != nil {
				return
			}
		default:
			err = fmt.Errorf("unexpexted token: %s: %s %d %d", tok.Type, tok.Text, p.indent, baseIndent)
			return
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

	spl := strings.Split(tok.Text, " ")
	switch len(spl) {
	case 1:
		field = &Field{Name: tok.Text}
	case 2:
		field = &Field{
			Name: spl[0],
			Type: spl[1],
		}
	}

	field.Description, err = p.readMultilineText(baseIndent + 1)
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

func (p *parser) readMultilineText(baseIndent int) (str string, err error) {
	for {
		tok := p.scan()
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

func (p *parser) readTextBlock(baseIndent int) (str string, err error) {
	for {
		tok := p.scan()
		if p.indent < baseIndent || tok.Type != TextTok {
			p.unscan()
			return
		}

		if str == "" {
			str = tok.Text
		} else {
			str += "\n" + tok.Text
		}
	}
}

func (p *parser) readExamples(baseIndent int) (egs []*Example, err error) {
	for {
		var eg *Example
		if eg, err = p.readExample(baseIndent + 1); err != nil || eg == nil {
			return egs, err
		}
		egs = append(egs, eg)
	}
}

func (p *parser) readExample(baseIndent int) (eg *Example, err error) {
	// read name
	tok := p.scan()
	if p.indent < baseIndent || tok.Type != TextTok {
		p.unscan()
		return nil, nil
	}

	eg = &Example{Name: tok.Text}
	for {
		tok := p.scan()
		if p.indent <= baseIndent {
			p.unscan()
			return
		}

		switch tok.Type {
		case CodeTok:
			if eg.Code, err = p.readTextBlock(p.indent); err != nil {
				return eg, err
			}
		case TextTok:
			p.unscan()
			if eg.Description, err = p.readTextBlock(p.indent); err != nil {
				return eg, err
			}
		default:
			p.unscan()
			return eg, nil
		}
	}
}

func (p *parser) errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}
