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
func (r *Runtime) Run(file string, nodes []ast.Node) error {
	_, err := r.Execute(file, false, nil, "global", r.with, nodes)
	return err
}

// Execute executes the given nodes and returns the result if any
func (r *Runtime) Execute(file string, snippet bool, parent Executer, scope string, pkgs map[string]packages.Package, nodes []ast.Node) ([]*packages.FuncReturn, error) {
	ex, execNodes, err := r.Executer(file, snippet, parent, scope, pkgs, nodes)
	if err != nil {
		return nil, err
	}
	return ex.Execute(execNodes)
}

func (r *Runtime) Executer(file string, snippet bool, parent Executer, scope string, pkgs map[string]packages.Package, nodes []ast.Node) (Executer, []ast.Node, error) {
	if snippet {
		ex := NewExecuter(r, parent, file, parent.GetNamespace(), scope, parent.GetPackages())
		return ex, nodes, nil
	}

	var (
		namespace = ""
		execNodes []ast.Node
	)

	for _, node := range nodes {
		switch node.Type {
		case ast.Namespace:
			if namespace != "" {
				return nil, nil, fmt.Errorf("namespace already defined")
			}
			namespace = node.Name
		case ast.UsePackage:
			if _, ok := pkgs[node.Name]; !ok {
				if _, ok := r.with[node.Name]; ok {
					pkgs[node.Value] = r.with[node.Name]
					continue
				}

				pkg := NewPackage(node.Name)
				if pkg == nil {
					return nil, nil, nodeErr(ErrPackageNotExists, node, fmt.Errorf("package '%s' not exists but used", node.Name))
				}
				pkgs[node.Value] = pkg
			}
		default:
			execNodes = append(execNodes, node)
		}
	}

	if namespace == "" {
		namespace = "main"
	}

	if parent != nil && namespace != parent.GetNamespace() {
		// if the namespace is different from the parent, then the parent is not the parent
		parent = nil
	}

	ex := NewExecuter(r, parent, file, namespace, scope, pkgs)
	return ex, execNodes, nil
}
