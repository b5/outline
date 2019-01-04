package lib

import (
	"bufio"
	"io"
	"strings"
)

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

type Scanner struct {
	r *bufio.Reader

	// scanning state
	tok               Token
	text              strings.Builder
	line, col, offset int
	readNewline       bool
	err               error
}

func (s *Scanner) Scan() Token {
	inText := false
	s.text = strings.Builder{}

	if s.readNewline {
		s.readNewline = false
		return s.newTok(NewlineTok)
	}

	for {
		ch := s.read()

		switch ch {
		case eof:
			if inText {
				s.readNewline = true
				return s.newTok(TextTok)
			}
			return s.newTok(EofTok)
		// ignore line feeds
		case '\r':
			continue
		case '\n':
			s.line++
			if inText {
				s.readNewline = true
				return s.newTok(TextTok)
			}
			return s.newTok(NewlineTok)
		case '\t':
			return s.newTok(IndentTok)
		case ':':
			switch s.text.String() {
			case "outline":
				return s.newTok(DocumentTok)
			case "functions":
				return s.newTok(FunctionsTok)
			case "methods":
				return s.newTok(MethodsTok)
			case "types":
				return s.newTok(TypesTok)
			case "fields":
				return s.newTok(FieldsTok)
			case "operators":
				return s.newTok(OperatorsTok)
			case "params":
				return s.newTok(ParamsTok)
			case "return":
				return s.newTok(ReturnTok)
			default:
				s.text.WriteRune(':')
			}
		case ' ':
			s.text.WriteRune(' ')
			if s.text.String() == "  " {
				return s.newTok(IndentTok)
			}
		default:
			s.text.WriteRune(ch)
			if !inText {
				inText = true
			}
		}
	}
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) newTok(t TokenType) Token {
	return Token{
		Type: t,
		Text: strings.TrimSpace(s.text.String()),
		Pos:  Position{Line: s.line, Col: s.col, Offset: s.offset},
	}
}

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
