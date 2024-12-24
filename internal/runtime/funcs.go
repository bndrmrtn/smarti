package runtime

import (
	"github.com/bndrmrtn/smarti/internal/ast"
)

type funcDecl struct {
	Args []ast.Node
	Body []ast.Node
}
