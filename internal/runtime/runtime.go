package runtime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/smlgh/smarti/internal/ast"
	"github.com/smlgh/smarti/internal/packages"
)

type Runtime struct {
	nodes []ast.Node

	uses      map[string]packages.Package
	variables map[string]variable

	with map[string]packages.Package

	mu sync.Mutex
}

func New(nodes []ast.Node) *Runtime {
	return &Runtime{
		nodes:     nodes,
		uses:      make(map[string]packages.Package),
		variables: make(map[string]variable),
		with:      make(map[string]packages.Package),
	}
}

func (r *Runtime) With(pkgName string, pkg packages.Package) {
	r.mu.Lock()
	r.with[pkgName] = pkg
	r.mu.Unlock()
}

func (r *Runtime) Run() error {
	for _, node := range r.nodes {
		switch node.Type {
		case ast.UsePackage:
			if _, ok := r.uses[node.Name]; !ok {
				if _, ok := r.with[node.Name]; ok {
					r.uses[node.Value] = r.with[node.Name]
					continue
				}

				pkg := NewPackage(node.Name)
				if pkg == nil {
					return fmt.Errorf("package %s not found", node.Name)
				}
				r.uses[node.Value] = pkg
			}
		case ast.VarExpression, ast.VarNil, ast.VarString, ast.VarSingleString, ast.VarNumber, ast.VarFloat, ast.VarBool, ast.VarTemplate, ast.VarVariable, ast.VarUnknown:
			if _, err := r.makeVar(node); err != nil {
				return err
			}
			break
		case ast.FuncCall:
			if _, err := r.callFn(node); err != nil {
				return err
			}
			break
		}
	}
	return nil
}

func (r *Runtime) makeVar(node ast.Node, ret ...bool) (interface{}, error) {
	var value interface{}

	switch node.Type {
	case ast.VarNil:
		value = nil
	case ast.VarString, ast.VarSingleString:
		value = node.Value
	case ast.VarNumber:
		v, err := strconv.Atoi(node.Value)
		if err != nil {
			return nil, err
		}
		value = v
		break
	case ast.VarFloat:
		v, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return nil, err
		}
		value = v
		break
	case ast.VarBool:
		v, err := strconv.ParseBool(node.Value)
		if err != nil {
			return nil, err
		}
		value = v
		break
	case ast.VarExpression:
		v, typ, err := r.handleExpression(node)
		if err != nil {
			return nil, err
		}
		node.Type = typ
		value = v
		break
	case ast.FuncCall:
		v, err := r.callFn(node)
		if err != nil {
			return nil, err
		}

		if len(v) == 0 {
			value = nil
			node.Type = ast.VarNil
			break
		}

		value = v[0].Value
		node.Type = toNodeType(v[0].Type)
		break
	case ast.VarTemplate:
		value = r.parseTemplate(node)
		node.Type = ast.VarString
		break
	}

	if len(ret) > 0 && ret[0] {
		return value, nil
	}

	r.mu.Lock()
	r.variables[node.Name] = variable{
		Type:  node.Type,
		Ref:   node.IsReference,
		Value: value,
	}
	r.mu.Unlock()

	return nil, nil
}

func (r *Runtime) callFn(node ast.Node) ([]packages.FuncReturn, error) {
	v, err := r.getArgs(node.Args)
	if err != nil {
		return nil, err
	}

	if strings.Contains(node.Name, ".") {
		parts := strings.Split(node.Name, ".")
		pkg, ok := r.uses[parts[0]]
		if !ok {
			return nil, fmt.Errorf("package %s not imported", parts[0])
		}
		return pkg.Run(parts[1], toPkgVar(v))
	}

	return builtin.runFn(node.Name, toPkgVar(v))
}

func (r *Runtime) getArgs(nodes []ast.Node) ([]variable, error) {
	args := make([]variable, len(nodes))
	for i, node := range nodes {
		switch node.Type {
		case ast.VarVariable:
			v, ok := r.variables[node.Value]
			if !ok {
				return nil, ast.NewErrWithPos(node.Info, errors.Join(ast.ErrorCannotUseBeforeDecl, fmt.Errorf("variable %s not declared", node.Value)))
			}
			args[i] = v
		default:
			v, err := r.makeVar(node, true)
			if err != nil {
				return nil, err
			}
			args[i] = variable{
				Type:  node.Type,
				Value: v,
			}
		}
	}
	return args, nil
}

func (r *Runtime) handleExpression(node ast.Node) (interface{}, ast.NodeType, error) {
	var (
		result interface{}
		typ    ast.NodeType
	)

	for _, n := range node.Children {
		switch n.Type {
		case ast.VarVariable:
			v, ok := r.variables[n.Value]
			if !ok {
				return nil, ast.VarUnknown, fmt.Errorf("variable %s not declared", n.Value)
			}
			result = v.Value
			typ = v.Type
			break
		case ast.VarExpression:
			v, t, err := r.handleExpression(n)
			if err != nil {
				return nil, ast.VarUnknown, err
			}
			result = v
			typ = t
			break
		case ast.FuncCall:
			v, err := r.callFn(n)
			if err != nil {
				return nil, ast.VarUnknown, err
			}
			if len(v) == 0 {
				return nil, ast.VarNil, nil
			}
			result = v[0].Value
			typ = toNodeType(v[0].Type)
			break
		default:
			result = nil
			typ = ast.VarNil
			break
		}
	}

	return result, typ, nil
}

func (r *Runtime) parseTemplate(node ast.Node) string {
	tmpl := parseTemplate(node.Value)
	fmt.Println(node.Value, tmpl)

	var tpl string

	for _, part := range tmpl {
		if part.Static {
			tpl += part.Content
			continue
		}

		val, typ, ref := ast.Type(part.Content)

		v, err := r.makeVar(ast.Node{
			Type:        typ,
			Value:       val,
			IsReference: ref,
		}, true)
		if err != nil {
			return ""
		}

		tpl += fmt.Sprintf("%v", v)
	}

	return tpl
}
