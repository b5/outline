package lib

type Position struct {
	Line, Col, Offset int
}

type Token struct {
	Type TokenType
	Pos  Position
	Text string
}

func (t Token) String() string {
	return t.Text
}

type TokenType int

const (
	IllegalTok TokenType = iota
	EofTok

	LiteralBegin
	IndentTok
	NewlineTok
	TextTok
	LiteralEnd

	KeywordBegin
	DocumentTok
	FunctionsTok
	TypesTok
	FieldsTok
	OperatorsTok
	KeywordEnd
)

func (t TokenType) String() string {
	switch t {
	case IndentTok:
		return "tab"
	case NewlineTok:
		return "newline"
	case TextTok:
		return "text"

	case DocumentTok:
		return "outline"
	case FunctionsTok:
		return "functions"
	case TypesTok:
		return "types"
	case FieldsTok:
		return "fields"
	case OperatorsTok:
		return "operators"
	default:
		return "unknown"
	}
}
