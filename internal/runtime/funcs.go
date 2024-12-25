package runtime

import (
	"github.com/bndrmrtn/smarti/internal/ast"
)

type funcDecl struct {
	Args []ast.Node
	Body []ast.Node
}

func getType(v any) ast.NodeType {
	switch v.(type) {
	case int:
		return ast.VarNumber
	case float64:
		return ast.VarFloat
	case string:
		return ast.VarString
	case bool:
		return ast.VarBool
	case nil:
		return ast.VarNil
	}
	return ast.VarUnknown
}
