package runtime

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/lexer"
	"github.com/bndrmrtn/smarti/internal/packages"
)

type CodeExecuter struct {
	file string

	parent    Executer
	scope     string
	namespace string
	runt      *Runtime

	uses      map[string]packages.Package
	variables map[string]*variable
	funcs     map[string]funcDecl

	children []Executer

	mu sync.Mutex
}

func NewExecuter(runt *Runtime, parent Executer, file, namespace, scope string, uses map[string]packages.Package) Executer {
	return &CodeExecuter{
		file:      filepath.Clean(file),
		parent:    parent,
		namespace: namespace,
		scope:     scope,
		uses:      uses,
		runt:      runt,
		variables: make(map[string]*variable),
		funcs:     make(map[string]funcDecl),
		children:  []Executer{},
	}
}

func (c *CodeExecuter) GetDir() string {
	return filepath.Dir(c.file)
}

func (c *CodeExecuter) GetNamespace() string {
	return c.namespace
}

func (c *CodeExecuter) GetParent() Executer {
	return c.parent
}

func (c *CodeExecuter) GetScope() string {
	return c.scope
}

func (c *CodeExecuter) runtime() *Runtime {
	return c.runt
}

func (c *CodeExecuter) GetPackage(name string) (packages.Package, error) {
	if pkg, ok := c.uses[name]; ok {
		return pkg, nil
	}

	return nil, ErrPackageNotImported
}

func (c *CodeExecuter) GetPackages() map[string]packages.Package {
	pkgs := c.uses
	if c.parent != nil {
		for k, v := range c.parent.GetPackages() {
			pkgs[k] = v
		}
	}
	return pkgs
}

func (c *CodeExecuter) DeclareVariable(name string, v *variable) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.variables[name]; ok {
		return ErrVariableAlreadyDeclared
	}

	c.variables[name] = v
	return nil
}

func (c *CodeExecuter) AssignVariable(name string, v *variable) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.variables[name]
	if !ok {
		if c.parent != nil {
			return c.parent.AssignVariable(name, v)
		}

		return ErrVariableNotDeclared
	}

	c.variables[name] = v
	return nil
}

func (c *CodeExecuter) DeclareFunc(name string, fn funcDecl) error {
	if c.parent != nil {
		return c.parent.DeclareFunc(name, fn)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.funcs[name]
	if ok {
		return fmt.Errorf("function %s() already declared", name)
	}

	c.funcs[name] = fn

	return nil
}

func (c *CodeExecuter) GetVariable(name string) (*variable, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.variables[name]
	if ok {
		return v, nil
	}

	if c.parent != nil {
		return c.parent.GetVariable(name)
	}

	return nil, ErrVariableNotDeclared
}

func (c *CodeExecuter) AccessVariableValue(name string) (*packages.Variable, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, err := c.GetVariable(name)
	if err != nil {
		return nil, err
	}

	if !v.Ref {
		return &packages.Variable{
			Type:  packages.VarType(v.Type),
			Value: v.Value,
		}, nil
	}

	v, err = c.GetVariable(v.Value.(string))
	if err != nil {
		return nil, err
	}

	return &packages.Variable{
		Type:  packages.VarType(v.Type),
		Value: v.Value,
	}, err
}

func (c *CodeExecuter) Execute(nodes []ast.Node) ([]*packages.FuncReturn, error) {
	for _, node := range nodes {
		switch node.Type {
		case ast.VarExpression, ast.VarNil, ast.VarString, ast.VarSingleString, ast.VarNumber, ast.VarFloat, ast.VarBool, ast.VarTemplate, ast.VarVariable, ast.VarUnknown:
			if _, _, err := c.createVariable(node); err != nil {
				return nil, err
			}
		case ast.FuncCall:
			if _, err := c.callFunc(node); err != nil {
				return nil, err
			}
		case ast.FuncDecl:
			c.DeclareFunc(node.Name, funcDecl{
				Args: node.Args,
				Body: node.Children,
			})
		case ast.FuncReturn:
			return c.funcGetReturn(node.Children)
		case ast.IfStatement:
			ok, err := c.evaluateStatement(node)
			if err != nil {
				return nil, err
			}
			if ok {
				ret, err := c.Execute(node.Children)
				if err != nil {
					return nil, err
				}

				if ret != nil {
					return ret, nil
				}
			}
			/*case ast.ForLoop:
			loopEx := NewExecuter(c.runt, c, c.file, c.namespace, "for", c.uses)
			initial := node.Args[0]*/
		}
	}

	if c.namespace == "main" && c.parent == nil && c.scope == "global" {
		if _, ok := c.funcs["main"]; !ok {
			return nil, nil
		}

		_, err := c.callFunc(ast.Node{
			Token: lexer.FuncCall,
			Name:  "main",
			Type:  ast.FuncCall,
		})
		return nil, err
	}

	return nil, nil
}
