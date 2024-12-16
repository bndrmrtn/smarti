package lexer

type ErrorPos struct {
	File string
	Line int
	Pos  int
}

type LexerErrorWithPosition struct {
	Pos *ErrorPos
	Err string
}

func NewLexerErrorWithPosition(pos *ErrorPos, err string) LexerErrorWithPosition {
	return LexerErrorWithPosition{
		Pos: pos,
		Err: err,
	}
}

func (l LexerErrorWithPosition) Error() string {
	return l.Err
}
