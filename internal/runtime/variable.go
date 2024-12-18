package runtime

import "github.com/smlgh/smarti/internal/ast"

type variable struct {
	Type  ast.NodeType
	Value interface{}
	Ref   bool
}
