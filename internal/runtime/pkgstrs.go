package runtime

import (
	"fmt"
	"strings"

	"github.com/smlgh/smarti/internal/ast"
)

type Strs struct{}

func (s Strs) Run(fn string, args []variable) ([]funcReturn, error) {
	switch fn {
	case "length":
		return s.fnLength(args)
	case "trim":
		return s.fnTrim(args)
	}
	return nil, nil
}

func (Strs) fnLength(args []variable) ([]funcReturn, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("length expects at least one argument")
	}
	arg := args[0]
	if arg.Type != ast.VarString && arg.Type != ast.VarSingleString {
		return nil, fmt.Errorf("length expects a string argument")
	}

	l := len(arg.Value.(string))
	return []funcReturn{
		{
			Value: l,
			Type:  ast.VarNumber,
		},
	}, nil
}

func (Strs) fnTrim(args []variable) ([]funcReturn, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("trim expects at least one argument")
	}
	arg := args[0]
	if arg.Type != ast.VarString && arg.Type != ast.VarSingleString {
		return nil, fmt.Errorf("trim expects a string argument")
	}

	return []funcReturn{
		{
			Value: strings.TrimSpace(arg.Value.(string)),
			Type:  ast.VarString,
		},
	}, nil
}
