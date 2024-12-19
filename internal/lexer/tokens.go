package lexer

type LexerToken struct {
	Type  Token  `json:"token"`
	Value string `json:"value"`
	Info  struct {
		File string `json:"file"`
		Line int    `json:"line"`
		Pos  int    `json:"pos"`
	}
}

func newLexerToken(t Token, v string, file string, line, pos int) LexerToken {
	return LexerToken{
		Type:  t,
		Value: v,
		Info: struct {
			File string `json:"file"`
			Line int    `json:"line"`
			Pos  int    `json:"pos"`
		}{
			File: file,
			Line: line + 1,
			Pos:  pos,
		},
	}
}

type Token int

const (
	// Use imports a package
	Use Token = iota

	// Let creates a new variable
	Let
	// Const creates a new constant
	Const
	// Equal assigns a value to a variable
	Assign

	// Nil is an uninitialized token
	Nil
	// StringLiteral is a string token: "
	DoubleStringLiteral
	// SingleStringLiteral is a single string token: '
	SingleStringLiteral

	// End marks the end of a block: semi-colon
	SemiColon

	// SingleLineComment is a single line comment: //
	SingleLineComment
	// MultiLineCommentStart is the start of a multi-line comment: /*
	MultiLineCommentStart
	// MultiLineCommentEnd is the end of a multi-line comment: */
	MultiLineCommentEnd

	// TemplateStart marks the start of an HTML block: <template>
	TemplateStart
	// TemplateEnd marks the end of an HTML block: </template>
	TemplateEnd
	Template

	// Export exports a variable
	Export Token = iota + 100
	// Use imports a package
	UseToken
	// Func creates a new function
	Func
	// FuncCall calls a function
	FuncCall
	// ParantesisStart is the start of a parantesis: (
	ParantesisStart
	// ParantesisEnd is the end of a parantesis: )
	ParantesisEnd
	// CurlyBraceStart is the start of a curly brace: {
	CurlyBraceStart
	// CurlyBraceEnd is the end of a curly brace: }
	CurlyBraceEnd
	// Return returns a value
	Return
	// Identifier is an identifier
	Identifier

	Addition = iota + 1000
	Subtraction
	Multiplication
	Division
	Modulo

	// Unknown is an unknown token
	Unknown = iota + 10000
)

func (t Token) String() string {
	switch t {
	case Use:
		return "use"
	case Let:
		return "let"
	case Const:
		return "const"
	case Assign:
		return "="
	case Nil:
		return "nil"
	case DoubleStringLiteral:
		return "\""
	case SingleStringLiteral:
		return "'"
	case SemiColon:
		return ";"
	case SingleLineComment:
		return "//"
	case MultiLineCommentStart:
		return "/*"
	case MultiLineCommentEnd:
		return "*/"
	case TemplateStart:
		return "<>"
	case TemplateEnd:
		return "</>"
	case Export:
		return "export"
	case UseToken:
		return "use"
	case Func:
		return "func"
	case FuncCall:
		return "funcCall"
	case ParantesisStart:
		return "("
	case ParantesisEnd:
		return ")"
	case CurlyBraceStart:
		return "{"
	case CurlyBraceEnd:
		return "}"
	case Return:
		return "return"
	default:
		return "unknown"
	}
}

func isToken(s string) Token {
	switch s {
	case "use":
		return Use
	case "let":
		return Let
	case "const":
		return Const
	case "=":
		return Assign
	case "+":
		return Addition
	case "-":
		return Subtraction
	case "*":
		return Multiplication
	case "/":
		return Division
	case "%":
		return Modulo
	case "nil":
		return Nil
	case ";":
		return SemiColon
	case "\"":
		return DoubleStringLiteral
	case "'":
		return SingleStringLiteral
	case "//":
		return SingleLineComment
	case "/*":
		return MultiLineCommentStart
	case "*/":
		return MultiLineCommentEnd
	case "<>":
		return TemplateStart
	case "</>":
		return TemplateEnd
	default:
		return Identifier
	}
}
