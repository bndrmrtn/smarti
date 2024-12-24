package runtime

import (
	"fmt"
	"sync"

	"github.com/bndrmrtn/smarti/internal/ast"

	"github.com/bndrmrtn/smarti/internal/packages"
)

type Runtime struct {
	with map[string]packages.Package

	mu sync.Mutex
}

func New() *Runtime {
	return &Runtime{
		with: make(map[string]packages.Package),
	}
}

func (r *Runtime) With(pkgName string, pkg packages.Package) {
	r.mu.Lock()
	r.with[pkgName] = pkg
	r.mu.Unlock()
}

// Run executes the given nodes as a main program
func (r *Runtime) Run(nodes []ast.Node) error {
	_, err := r.Execute(nil, "global", r.with, nodes)
	return err
}

// Execute executes the given nodes and returns the result if any
func (r *Runtime) Execute(parent Executer, scope string, pkgs map[string]packages.Package, nodes []ast.Node) ([]*packages.FuncReturn, error) {
	var (
		namespace = "main"
		execNodes []ast.Node
	)

	for _, node := range nodes {
		switch node.Type {
		case ast.Namespace:
			namespace = node.Name
		case ast.UsePackage:
			if _, ok := pkgs[node.Name]; !ok {
				if _, ok := r.with[node.Name]; ok {
					pkgs[node.Value] = r.with[node.Name]
					continue
				}

				pkg := NewPackage(node.Name)
				if pkg == nil {
					return nil, fmt.Errorf("package %s not found", node.Name)
				}
				pkgs[node.Value] = pkg
			}
		default:
			execNodes = append(execNodes, node)
		}
	}

	ex := NewExecuter(r, parent, namespace, scope, pkgs)
	return ex.Execute(execNodes)
}
