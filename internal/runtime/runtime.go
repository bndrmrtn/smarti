package runtime

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/smlgh/smarti/internal/ast"
)

type Runtime struct {
	nodes []ast.Node

	variables map[string]variable

	mu sync.Mutex
}

func New(nodes []ast.Node) *Runtime {
	return &Runtime{
		nodes:     nodes,
		variables: make(map[string]variable),
	}
}

func (r *Runtime) Run() error {
	for _, node := range r.nodes {
		switch node.Type {
		case ast.VarNil, ast.VarString, ast.VarSingleString, ast.VarNumber, ast.VarFloat, ast.VarBool, ast.VarTemplate, ast.VarVariable, ast.VarUnknown:
			r.makeVar(node)
		case ast.FuncCall:
			v, err := r.getArgs(node.Args)
			if err != nil {
				return err
			}
			builtin.runFn(node.Name, v)
			break
		}
	}
	return nil
}

func (r *Runtime) makeVar(node ast.Node, ret ...bool) interface{} {
	var value interface{}

	switch node.Type {
	case ast.VarNil:
		value = nil
	case ast.VarString, ast.VarSingleString:
		value = node.Value
	case ast.VarNumber:
		v, _ := strconv.Atoi(node.Value)
		value = v
		break
	case ast.VarFloat:
		v, _ := strconv.ParseFloat(node.Value, 64)
		value = v
		break
	case ast.VarBool:
		v, _ := strconv.ParseBool(node.Value)
		value = v
		break
	case ast.VarTemplate:
		value = node.Value
	}

	if len(ret) > 0 && ret[0] {
		return value
	}

	r.mu.Lock()
	r.variables[node.Name] = variable{
		Type:  node.Type,
		Ref:   node.IsReference,
		Value: value,
	}
	r.mu.Unlock()

	return nil
}

func (r *Runtime) getArgs(nodes []ast.Node) ([]variable, error) {
	args := make([]variable, len(nodes))
	for i, node := range nodes {
		switch node.Type {
		case ast.VarVariable:
			v, ok := r.variables[node.Value]
			if !ok {
				fmt.Println(node)
				return nil, ast.NewErrWithPos(node.Info, errors.Join(ast.ErrorCannotUseBeforeDecl, fmt.Errorf("variable %s not declared", node.Value)))
			}
			args[i] = v
		default:
			v := r.makeVar(node, true)
			args[i] = variable{
				Type:  node.Type,
				Value: v,
			}
		}
	}
	return args, nil
}
