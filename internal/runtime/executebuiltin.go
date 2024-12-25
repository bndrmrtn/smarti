package runtime

import (
	"fmt"
	"path/filepath"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/lexer"
	"github.com/bndrmrtn/smarti/internal/packages"
)

func (c *CodeExecuter) ExecuteBuiltinMethod(e Executer, name string, args []*packages.Variable) ([]*packages.FuncReturn, error) {
	switch name {
	case "type":
		return runFnType(args)
	case "import":
		return runFnImport(e, args)
	}
	return nil, fmt.Errorf("function %s does not exists or imported", name)
}

func runFnType(args []*packages.Variable) ([]*packages.FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type function expects 1 argument, %d given", len(args))
	}

	return []*packages.FuncReturn{
		{
			Value: args[0].Type,
			Type:  packages.VarString,
		},
	}, nil
}

func runFnImport(e Executer, args []*packages.Variable) ([]*packages.FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("import function expects 1 argument, %d given", len(args))
	}

	if args[0].Type != packages.VarString && args[0].Type != packages.VarSingleString {
		return nil, fmt.Errorf("import function expects string argument, %s given", args[0].Type)
	}

	file := filepath.Join(e.GetDir(), args[0].Value.(string))
	lx := lexer.New(file)
	if err := lx.Parse(); err != nil {
		return nil, err
	}

	ps := ast.NewParser(lx.Tokens)
	if err := ps.Parse(); err != nil {
		return nil, err
	}

	return e.Execute(ps.Nodes)
}
