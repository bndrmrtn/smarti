package runtime

import (
	"fmt"

	"github.com/bndrmrtn/smarti/internal/ast"
)

type Err error

var (
	ErrVariableNotDeclared     Err = fmt.Errorf("variable not declared")
	ErrFuncNotDeclared         Err = fmt.Errorf("function not declared")
	ErrVariableAlreadyDeclared Err = fmt.Errorf("variable already declared")
	ErrPackageNotImported      Err = fmt.Errorf("package not imported")
	ErrNotExpression           Err = fmt.Errorf("not an expression")
	ErrInvalidExpression       Err = fmt.Errorf("invalid expression")
	ErrFuncCall                Err = fmt.Errorf("function call error")
	ErrVariable                Err = fmt.Errorf("variable error")
)

func nodeErr(typ Err, n ast.Node, err error) error {
	return fmt.Errorf("%w: %w at %s", typ, err, n.Info.String())
}
