package runtime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/packages"
)

type Runtime struct {
	uses      map[string]packages.Package
	variables map[string]variable

	funcs map[string]funcDecl

	with map[string]packages.Package

	mu sync.Mutex
}

var builtin runtimeBuiltin

func New() *Runtime {
	r := &Runtime{
		uses:      make(map[string]packages.Package),
		variables: make(map[string]variable),
		funcs:     make(map[string]funcDecl),
		with:      make(map[string]packages.Package),
	}
	builtin = runtimeBuiltin{r: r}
	return r
}

func (r *Runtime) With(pkgName string, pkg packages.Package) {
	r.mu.Lock()
	r.with[pkgName] = pkg
	r.mu.Unlock()
}

func (r *Runtime) Run(nodes []ast.Node) error {
	if _, err := r.Execute(nodes); err != nil {
		return err
	}

	main, ok := r.funcs["main"]
	if !ok {
		return nil
	}

	_, err := r.Execute(main.Body)
	return err
}

func (r *Runtime) Execute(nodes []ast.Node) ([]packages.FuncReturn, error) {
	var scoped []ast.Node

	defer func() {
		for _, node := range scoped {
			switch node.Type {
			case ast.VarExpression, ast.VarNil, ast.VarString, ast.VarSingleString, ast.VarNumber, ast.VarFloat, ast.VarBool, ast.VarTemplate, ast.VarVariable, ast.VarUnknown:
				delete(r.variables, node.Name)
			}
		}
	}()

	for _, node := range nodes {
		if node.Scope != ast.ScopeGlobal {
			scoped = append(scoped, node)
		}

		switch node.Type {
		case ast.UsePackage:
			if _, ok := r.uses[node.Name]; !ok {
				if _, ok := r.with[node.Name]; ok {
					r.uses[node.Value] = r.with[node.Name]
					continue
				}

				pkg := NewPackage(node.Name)
				if pkg == nil {
					return nil, fmt.Errorf("package %s not found", node.Name)
				}
				r.uses[node.Value] = pkg
			}
		case ast.VarExpression, ast.VarNil, ast.VarString, ast.VarSingleString, ast.VarNumber, ast.VarFloat, ast.VarBool, ast.VarTemplate, ast.VarVariable, ast.VarUnknown:
			if _, err := r.makeVar(node); err != nil {
				return nil, err
			}
			break
		case ast.FuncCall:
			if _, err := r.callFn(node); err != nil {
				return nil, err
			}
			break
		case ast.FuncDecl:
			r.funcs[node.Name] = funcDecl{
				Args: node.Args,
				Body: node.Children,
			}
		case ast.FuncReturn:
			return r.makeReturn(node.Children)
		}
	}

	return nil, nil
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
	case ast.VarVariable:
		v, ok := r.variables[node.Value]
		if !ok {
			return nil, fmt.Errorf("variable %s not declared", node.Value)
		}
		value = v.Value
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
		return pkg.Run(parts[1], r.toPkgVar(v))
	}

	fn, ok := r.funcs[node.Name]
	if ok {
		return r.Execute(fn.Body)
	}

	return builtin.runFn(node.Name, r.toPkgVar(v))
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

func (r *Runtime) makeReturn(nodes []ast.Node) ([]packages.FuncReturn, error) {
	var returns []packages.FuncReturn

	for _, v := range nodes {
		val, err := r.makeVar(v, true)
		if err != nil {
			return nil, err
		}

		if v.Type == ast.VarVariable {
			vv, ok := r.variables[v.Value]
			if !ok {
				return nil, fmt.Errorf("variable %s not declared", v.Value)
			}
			v.Type = vv.Type
		}

		returns = append(returns, packages.FuncReturn{
			Type:  toPkgType(v.Type),
			Value: val,
		})
	}

	return returns, nil
}
