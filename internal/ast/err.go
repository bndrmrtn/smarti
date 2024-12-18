package ast

import "github.com/fatih/color"

type ErrWithPos struct {
	Pos NodeFileInfo
	Err string
}

func NewErrWithPos(pos NodeFileInfo, err error) ErrWithPos {
	return ErrWithPos{
		Pos: pos,
		Err: err.Error(),
	}
}

func (l ErrWithPos) Error() string {
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	return red("Error: ") + l.Err + " at " + l.Pos.String()
}

type Err string

const (
	ErrorCannotReDeclareVar  Err = "cannot redeclare variable"
	ErrorCannotReAssignConst Err = "cannot reassign constant"
	ErrorCannotUseBeforeDecl Err = "cannot use variable before declaration"
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
