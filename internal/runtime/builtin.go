package runtime

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/lexer"
	"github.com/bndrmrtn/smarti/internal/packages"
)

type runtimeBuiltin struct {
	r *Runtime
}

type funcReturn struct {
	Value interface{}
	Type  ast.NodeType
}

func (b *runtimeBuiltin) runFn(fn string, args []packages.Variable) ([]packages.FuncReturn, error) {
	switch fn {
	case "import":
		return builtin.runFnImport(args)
	case "type":
		return builtin.runFnType(args)
	case "writeType":
		return b.runFnWriteType(args)
	case "capitalize":
		return b.runFnCapitalize(args)
	case "format":
		return b.runFnFormat(args)
	}
	return nil, fmt.Errorf("func \"%s\" does not exists", fn)
}

func (b *runtimeBuiltin) runFnType(args []packages.Variable) ([]packages.FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type function expects exactly one argument")
	}
	return []packages.FuncReturn{
		{
			Value: string(args[0].Type),
			Type:  packages.VarString,
		},
	}, nil
}

func (b *runtimeBuiltin) runFnWriteType(args []packages.Variable) ([]packages.FuncReturn, error) {
	var types []interface{}
	for _, arg := range args {
		types = append(types, string(arg.Type))
	}
	fmt.Print(types...)
	return nil, nil
}

func (b *runtimeBuiltin) runFnCapitalize(args []packages.Variable) ([]packages.FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("capitalize function expects exactly one argument")
	}
	if args[0].Type != packages.VarString {
		return nil, fmt.Errorf("capitalize function expects string as argument")
	}
	s := args[0].Value.(string)
	return []packages.FuncReturn{
		{
			Type:  packages.VarString,
			Value: strings.ToUpper(s[:1]) + s[1:],
		},
	}, nil
}

func (b *runtimeBuiltin) runFnFormat(args []packages.Variable) ([]packages.FuncReturn, error) {
	var format string
	values := make([]interface{}, len(args)-1)
	for i, arg := range args {
		if i == 0 {
			if arg.Type == packages.VarString || arg.Type == packages.VarSingleString {
				format = arg.Value.(string)
			} else {
				return nil, fmt.Errorf("writef expects first argument to be a string")
			}
			continue
		}
		values[i-1] = arg.Value
	}
	val := fmt.Sprintf(format, values...)
	return []packages.FuncReturn{
		{
			Type:  packages.VarString,
			Value: val,
		},
	}, nil
}

func (b *runtimeBuiltin) runFnImport(args []packages.Variable) ([]packages.FuncReturn, error) {
	var main string
	var files []string
	for i, arg := range args {
		if arg.Type != packages.VarString {
			return nil, fmt.Errorf("import function expects string as argument")
		}

		if i == 0 {
			main = arg.Value.(string)
			continue
		}

		files = append(files, arg.Value.(string))
	}

	lx := lexer.New(main, files...)
	if err := lx.Parse(); err != nil {
		return nil, err
	}

	a := ast.NewParser(lx.Tokens)
	if err := a.Parse(); err != nil {
		return nil, err
	}

	return b.r.Execute(a.Nodes)
}
