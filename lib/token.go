package lib

// Position of a token within the scan stream
type Position struct {
	Line, Col, Offset int
}

// Token is a recognized token from the outlineline lexicon
type Token struct {
	Type TokenType
	Pos  Position
	Text string
}

// String implements the stringer interface for token
func (t Token) String() string {
	return t.Text
}

// TokenType enumerates the different types of tokens
type TokenType int

const (
	// IllegalTok is the default for unrecognized tokens
	IllegalTok TokenType = iota
	eofTok

	// LiteralBegin marks the beginning of literal tokens in the token enumeration
	LiteralBegin
	// IndentTok is a tab character "\t" or two consecutive spaces"  "
	IndentTok
	// NewlineTok is a line break
	NewlineTok
	// TextTok is a token for arbitrary text
	TextTok
	// LiteralEnd marks the end of literal tokens in the token enumeration
	LiteralEnd

	// KeywordBegin marks the end of keyword tokens in the token enumeration
	KeywordBegin
	// CodeTok is the "code:" token
	CodeTok
	// DocumentTok is the "document:" token
	DocumentTok
	// ExamplesTok is the "examples:" token
	ExamplesTok
	// PathTok is the "path:" token
	PathTok
	// FunctionsTok is the "functions:" token
	FunctionsTok
	// ParamsTok is the "params:" token
	ParamsTok
	// ReturnTok is the "return:" token
	ReturnTok
	// TypesTok is the "types:" token
	TypesTok
	// FieldsTok is the "fields:" token
	FieldsTok
	// MethodsTok is the "methods:" token
	MethodsTok
	// OperatorsTok is the "operators:" token
	OperatorsTok
	// KeywordEnd marks the end of keyword tokens in the token enumeration
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
	case PathTok:
		return "path"
	case MethodsTok:
		return "methods"
	case ExamplesTok:
		return "examples"
	case CodeTok:
		return "code"
	case FunctionsTok:
		return "functions"
	case TypesTok:
		return "types"
	case FieldsTok:
		return "fields"
	case OperatorsTok:
		return "operators"
	case ParamsTok:
		return "params"
	case ReturnTok:
		return "return"
	default:
		return "unknown"
	}
}
