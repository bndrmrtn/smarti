package runtime

import (
	"github.com/bndrmrtn/smarti/internal/ast"

	"github.com/bndrmrtn/smarti/internal/packages"
)

type Executer interface {
	GetDir() string

	GetNamespace() string
	GetParent() Executer
	GetScope() string

	DeclareVariable(name string, v *variable) error
	AssignVariable(name string, v *variable) error
	GetVariable(name string) (*variable, error)
	AccessVariableValue(name string) (*packages.Variable, error)
	DeclareFunc(name string, fn funcDecl) error

	GetPackages() map[string]packages.Package
	GetPackage(name string) (packages.Package, error)

	Execute(nodes []ast.Node) ([]*packages.FuncReturn, error)

	// Core methods

	createVariable(node ast.Node, onlyReturnValue ...bool) (interface{}, ast.NodeType, error)
	callFunc(node ast.Node) ([]*packages.FuncReturn, error)
	funcGetArgs(nodes []ast.Node) ([]*variable, error)
	funcGetReturn(nodes []ast.Node) ([]*packages.FuncReturn, error)

	evaluateExpression(node ast.Node) (interface{}, ast.NodeType, error)
	evaluateTemplate(node ast.Node) (string, error)
	runtime() *Runtime
}
