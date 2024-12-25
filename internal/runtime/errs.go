package runtime

import (
	"errors"
	"fmt"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/fatih/color"
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
	ErrInvalidFuncArgument     Err = fmt.Errorf("invalid function argument")
	ErrInvalidFuncReturn       Err = fmt.Errorf("invalid function return")
)

func nodeErr(typ Err, n ast.Node, err error) error {
	var runtimeErr Err
	if errors.Is(err, runtimeErr) {
		return err
	}

	redB := color.New(color.FgRed, color.Bold).SprintfFunc()
	red := color.New(color.FgRed).SprintfFunc()
	yel := color.New(color.FgYellow).SprintfFunc()

	at := "unknown"
	if n.Info.File != "" {
		at = n.Info.String()
	}

	return fmt.Errorf("%s %s\n%s\n%s\n", redB("Error type:"), red("%v,", typ), yel("%v at:", err), at)
}
