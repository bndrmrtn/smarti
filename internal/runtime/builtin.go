package runtime

import (
	"fmt"
	"strings"

	"github.com/smlgh/smarti/internal/ast"
)

type runtimeBuiltin struct{}

type funcReturn struct {
	Value interface{}
	Type  ast.NodeType
}

var builtin runtimeBuiltin

func (b runtimeBuiltin) runFn(fn string, args []variable) ([]funcReturn, error) {
	switch fn {
	case "type":
		return builtin.runFnType(args)
	case "writeType":
		return b.runFnWriteType(args)
	case "capitalize":
		return b.runFnCapitalize(args)
	}
	return nil, fmt.Errorf("func \"%s\" does not exists", fn)
}

func (b runtimeBuiltin) runFnType(args []variable) ([]funcReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type function expects exactly one argument")
	}
	return []funcReturn{
		{
			Value: string(args[0].Type),
			Type:  ast.VarString,
		},
	}, nil
}

func (b runtimeBuiltin) runFnWriteType(args []variable) ([]funcReturn, error) {
	var types []interface{}
	for _, arg := range args {
		types = append(types, string(arg.Type))
	}
	fmt.Print(types...)
	return nil, nil
}

func (b runtimeBuiltin) runFnCapitalize(args []variable) ([]funcReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("capitalize function expects exactly one argument")
	}
	if args[0].Type != ast.VarString {
		return nil, fmt.Errorf("capitalize function expects string as argument")
	}
	s := args[0].Value.(string)
	return []funcReturn{
		{
			Type:  ast.VarString,
			Value: strings.ToUpper(s[:1]) + s[1:],
		},
	}, nil
}
