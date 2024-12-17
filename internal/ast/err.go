package ast

type ErrWithPos struct {
	Pos NodeFileInfo
	Err string
}

func NewErrWithPos(pos NodeFileInfo, err string) ErrWithPos {
	return ErrWithPos{
		Pos: pos,
		Err: err,
	}
}

func (l ErrWithPos) Error() string {
	return l.Err
}

type Err string

const (
	ErrorCannotReDeclareVar  Err = "cannot redeclare variable"
	ErrorCannotReAssignConst Err = "cannot reassign constant"
	ErrorInvalidToken        Err = "invalid token"
	ErrorUnexpectedToken     Err = "unexpected token"
	ErrorUnexpectedEOF       Err = "unexpected EOF"
	ErrorInvalidSyntax       Err = "invalid syntax"
	ErrorInvalidType         Err = "invalid type"
	ErrorInvalidValue        Err = "invalid value"
	ErrorInvalidOperator     Err = "invalid operator"
	ErrorInvalidExpression   Err = "invalid expression"
	ErrorInvalidStatement    Err = "invalid statement"
	ErrorInvalidFunction     Err = "invalid function"
	ErrorInvalidParameter    Err = "invalid parameter"
	ErrorInvalidReturn       Err = "invalid return"
	ErrorInvalidAssignment   Err = "invalid assignment"
	ErrorInvalidCondition    Err = "invalid condition"
	ErrorInvalidLoop         Err = "invalid loop"
	ErrorInvalidBreak        Err = "invalid break"
	ErrorInvalidContinue     Err = "invalid continue"
	ErrorInvalidBlock        Err = "invalid block"
	ErrorInvalidCall         Err = "invalid call"
	ErrorInvalidIndex        Err = "invalid index"
	ErrorInvalidField        Err = "invalid field"
	ErrorInvalidMethod       Err = "invalid method"
)

func (e Err) Error() string {
	return string(e)
}
