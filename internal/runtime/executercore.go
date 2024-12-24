package runtime

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bndrmrtn/smarti/internal/ast"
	"github.com/bndrmrtn/smarti/internal/packages"
)

func (c *CodeExecuter) createVariable(node ast.Node, ret ...bool) (interface{}, ast.NodeType, error) {
	var value interface{}

	switch node.Type {
	case ast.VarNil:
		value = nil
	case ast.VarString, ast.VarSingleString:
		value = node.Value
	case ast.VarNumber:
		v, err := strconv.Atoi(node.Value)
		if err != nil {
			return nil, ast.VarUnknown, nodeErr(ErrVariable, node, fmt.Errorf("invalid number %v", node.Value))
		}
		value = v
		break
	case ast.VarFloat:
		v, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return nil, ast.VarUnknown, nodeErr(ErrVariable, node, fmt.Errorf("invalid float %v", node.Value))
		}
		value = v
		break
	case ast.VarBool:
		v, err := strconv.ParseBool(node.Value)
		if err != nil {
			return nil, ast.VarUnknown, err
		}
		value = v
		break
	case ast.VarExpression:
		v, typ, err := c.evaluateExpression(node)
		if err != nil {
			return nil, ast.VarUnknown, err
		}
		node.Type = typ
		value = v
		break
	case ast.FuncCall:
		v, err := c.callFunc(node)
		if err != nil {
			return nil, ast.VarUnknown, err
		}

		if len(v) == 0 {
			value = nil
			node.Type = ast.VarNil
			break
		}

		value = v[0].Value
		node.Type = toNodeType(v[0].Type)
		break
	case ast.VarVariable:
		v, err := c.GetVariable(node.Value)
		if err != nil {
			return nil, ast.VarUnknown, err
		}
		node.Type = v.Type
		value = v.Value
	}

	if len(ret) > 0 && ret[0] {
		return value, node.Type, nil
	}

	c.mu.Lock()
	c.variables[node.Name] = &variable{
		Type:  node.Type,
		Ref:   node.IsReference,
		Value: value,
	}
	c.mu.Unlock()

	return nil, node.Type, nil
}

func (c *CodeExecuter) callFunc(node ast.Node) ([]*packages.FuncReturn, error) {
	v, err := c.funcGetArgs(node.Args)
	if err != nil {
		return nil, err
	}

	if strings.Contains(node.Name, ".") {
		parts := strings.Split(node.Name, ".")
		pkg, ok := c.uses[parts[0]]
		if !ok {
			return nil, fmt.Errorf("package %s not imported", parts[0])
		}
		return pkg.Run(parts[1], toPkgVar(v))
	}

	fn, ok := c.funcs[node.Name]
	if ok {
		return c.runt.Execute(c, "global", c.GetPackages(), fn.Body)
	}

	// TODO: Check if it is a builtin function
	return c.ExecuteBuiltinMethod(c, node.Name, toPkgVar(v))
}

func (c *CodeExecuter) funcGetArgs(nodes []ast.Node) ([]*variable, error) {
	args := make([]*variable, len(nodes))
	for i, node := range nodes {
		switch node.Type {
		case ast.VarVariable:
			v, err := c.GetVariable(node.Name)
			if err != nil {
				return nil, err
			}
			args[i] = v
		default:
			v, t, err := c.createVariable(node, true)
			if err != nil {
				return nil, err
			}
			args[i] = &variable{
				Type:  t,
				Value: v,
			}
		}
	}
	return args, nil
}

func (c *CodeExecuter) funcGetReturn(nodes []ast.Node) ([]*packages.FuncReturn, error) {
	var returns []*packages.FuncReturn

	for _, v := range nodes {
		val, t, err := c.createVariable(v, true)
		if err != nil {
			return nil, err
		}

		returns = append(returns, &packages.FuncReturn{
			Type:  toPkgType(t),
			Value: val,
		})
	}

	return returns, nil
}

func (c *CodeExecuter) evaluateExpression(node ast.Node) (interface{}, ast.NodeType, error) {
	if node.Type != ast.VarExpression {
		return nil, ast.VarNil, ErrNotExpression
	}

	type expr struct {
		Op    bool
		Value interface{}
	}

	var (
		result []expr
		typ    ast.NodeType
	)

	for _, n := range node.Children {
		switch n.Type {
		case ast.VarOperator:
			if typ == "" {
				return nil, ast.VarUnknown, ErrInvalidExpression
			}

			if len(result) > 1 && result[len(result)-1].Op {
				return nil, ast.VarUnknown, ErrInvalidExpression
			}

			result = append(result, expr{
				Op:    true,
				Value: n.Value,
			})
		default:
			v, t, err := c.createVariable(n, true)
			if err != nil {
				return nil, ast.VarUnknown, err
			}

			if typ == "" {
				typ = t
			}

			if typ != t {
				return nil, ast.VarUnknown, ErrInvalidExpression
			}

			result = append(result, expr{
				Op:    false,
				Value: v,
			})
		}
	}

	if len(result) == 1 {
		return result[0].Value, typ, nil
	}

	var out interface{}
	var i = 0
	for i+3 <= len(result) {
		out = result[i].Value
		switch typ {
		case ast.VarNumber:
			if v, ok := result[i].Value.(int); ok {
				switch result[i+1].Value {
				case "+":
					out = v + result[i+2].Value.(int)
				case "-":
					out = v - result[i+2].Value.(int)
				case "*":
					out = v * result[i+2].Value.(int)
				case "/":
					out = v / result[i+2].Value.(int)
				}
			} else {
				return nil, ast.VarUnknown, ErrInvalidExpression
			}
		case ast.VarString, ast.VarSingleString:
			if v, ok := result[i].Value.(string); ok {
				switch result[i+1].Value {
				case "+":
					out = v + result[i+2].Value.(string)
				case "-":
					out = strings.Replace(v, result[i+2].Value.(string), "", -1)
				}
			} else {
				return nil, ast.VarUnknown, ErrInvalidExpression
			}
		}
		i += 3
	}

	return out, typ, nil
}
