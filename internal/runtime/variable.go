package runtime

import (
	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/packages"
)

type variable struct {
	Type  ast.NodeType
	Value interface{}
	Ref   bool
	Scope ast.NodeScope
}

func (r *Runtime) toPkgVar(v []variable) []packages.Variable {
	var vars []packages.Variable

	for _, vv := range v {
		vars = append(vars, packages.Variable{
			Type:  packages.VarType(vv.Type),
			Value: vv.Value,
		})
	}

	return vars
}

func toNodeType(v packages.VarType) ast.NodeType {
	return ast.NodeType(v)
}

func toPkgType(t ast.NodeType) packages.VarType {
	return packages.VarType(t)
}
