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
			return nil, ast.VarUnknown, nodeErr(ErrVariable, node, fmt.Errorf("invalid number: %v", node.Value))
		}
		value = v
		break
	case ast.VarFloat:
		v, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return nil, ast.VarUnknown, nodeErr(ErrVariable, node, fmt.Errorf("invalid float: %v", node.Value))
		}
		value = v
		break
	case ast.VarBool:
		v, err := strconv.ParseBool(node.Value)
		if err != nil {
			return nil, ast.VarUnknown, nodeErr(ErrVariable, node, fmt.Errorf("invalid boolean: %v", node.Value))
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
			return nil, ast.VarUnknown, nodeErr(ErrVariable, node, fmt.Errorf("invalid variable reference: '%v'", node.Value))
		}
		node.Type = v.Type
		value = v.Value
	case ast.VarTemplate:
		v, err := c.evaluateTemplate(node)
		if err != nil {
			return nil, ast.VarUnknown, err
		}
		value = v
		node.Type = ast.VarString
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
			return nil, nodeErr(ErrPackageNotImported, node, fmt.Errorf("package %s not imported", parts[0]))
		}
		return pkg.Run(parts[1], toPkgVar(v))
	}

	fn, ok := c.funcs[node.Name]
	if ok {
		ex, nodes, err := c.runt.Executer(c.file, true, c, "func", c.GetPackages(), fn.Body)
		if err != nil {
			return nil, nodeErr(ErrFuncCall, node, err)
		}

		if len(fn.Args) != len(v) {
			return nil, nodeErr(ErrFuncCall, node, fmt.Errorf("invalid number of arguments. expected %d, got %d", len(fn.Args), len(v)))
		}

		for i, arg := range fn.Args {
			err := ex.DeclareVariable(arg.Value, v[i])
			if err != nil {
				return nil, nodeErr(ErrFuncCall, node, err)
			}
		}

		return ex.Execute(nodes)
	}

	if c.parent != nil {
		return c.parent.callFunc(node)
	}

	return c.ExecuteBuiltinMethod(c, node.Name, toPkgVar(v))
}

func (c *CodeExecuter) funcGetArgs(nodes []ast.Node) ([]*variable, error) {
	args := make([]*variable, len(nodes))
	for i, node := range nodes {
		switch node.Type {
		case ast.VarVariable:
			v, err := c.GetVariable(node.Value)
			if err != nil {
				return nil, nodeErr(ErrInvalidFuncArgument, node, fmt.Errorf("invalid variable reference: '%v'", node.Value))
			}
			args[i] = v
		default:
			v, t, err := c.createVariable(node, true)
			if err != nil {
				return nil, nodeErr(ErrInvalidFuncArgument, node, err)
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
			return nil, nodeErr(ErrInvalidFuncReturn, v, err)
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
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("operator cannot come before value"))
			}

			if len(result) > 0 && result[len(result)-1].Op {
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("two operators cannot follow each other"))
			}

			result = append(result, expr{
				Op:    true,
				Value: n.Value,
			})
		default:
			v, t, err := c.createVariable(n, true)
			if err != nil {
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, err)
			}

			if typ == "" {
				typ = t
			}

			if typ != t {
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("mismatched types in expression: %s vs %s", typ, t))
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

	var (
		compare bool
		out     interface{}
		i       int
	)

	for i+2 < len(result) {
		if result[i].Op || !result[i+1].Op || result[i+2].Op {
			return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("invalid expression structure"))
		}

		left, right := result[i].Value, result[i+2].Value

		switch typ {
		case ast.VarNumber:
			leftVal, leftOk := left.(int)
			rightVal, rightOk := right.(int)
			if !leftOk || !rightOk {
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("invalid number operands"))
			}

			switch result[i+1].Value {
			case "+":
				out = leftVal + rightVal
			case "-":
				out = leftVal - rightVal
			case "*":
				out = leftVal * rightVal
			case "/":
				if rightVal == 0 {
					return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("division by zero"))
				}
				out = leftVal / rightVal
			case "==":
				out = leftVal == rightVal
				compare = true
			case "!=":
				out = leftVal != rightVal
				compare = true
			default:
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("unsupported operator: %v", result[i+1].Value))
			}

		case ast.VarString, ast.VarSingleString:
			leftVal, leftOk := left.(string)
			rightVal, rightOk := right.(string)
			if !leftOk || !rightOk {
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("invalid string operands"))
			}

			switch result[i+1].Value {
			case "+":
				out = leftVal + rightVal
			case "-":
				out = strings.ReplaceAll(leftVal, rightVal, "")
			case "==":
				out = leftVal == rightVal
				compare = true
			case "!=":
				out = leftVal != rightVal
				compare = true
			default:
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("unsupported operator: %v", result[i+1].Value))
			}

		case ast.VarBool:
			leftVal, leftOk := left.(bool)
			rightVal, rightOk := right.(bool)
			if !leftOk || !rightOk {
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("invalid boolean operands"))
			}

			switch result[i+1].Value {
			case "==":
				out = leftVal == rightVal
				compare = true
			case "!=":
				out = leftVal != rightVal
				compare = true
			default:
				return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("unsupported operator for booleans: %v", result[i+1].Value))
			}

		default:
			return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("unsupported type: %v", typ))
		}

		result = append(result[:i], expr{Op: false, Value: out})
		result = append(result, result[i+1:]...)
	}

	if compare {
		typ = ast.VarBool
	}

	if len(result) == 1 {
		return result[0].Value, typ, nil
	}

	return nil, ast.VarUnknown, nodeErr(ErrInvalidExpression, node, fmt.Errorf("incomplete expression"))
}

func (c *CodeExecuter) evaluateTemplate(node ast.Node) (string, error) {
	parts := parseTemplate(node.Value)
	var sb strings.Builder

	for _, part := range parts {
		if part.Static {
			sb.WriteString(part.Content)
			continue
		}
		v, err := c.GetVariable(part.Content)
		if err != nil {
			return "", nodeErr(ErrInvalidTemplate, node, err)
		}
		sb.WriteString(fmt.Sprint(v.Value))
	}

	return sb.String(), nil
}

func (c *CodeExecuter) evaluateStatement(node ast.Node) (bool, error) {
	ok, typ, err := c.evaluateExpression(ast.Node{
		Type:     ast.VarExpression,
		Children: node.Args,
	})
	if err != nil {
		return false, err
	}

	if typ != ast.VarBool {
		return false, nodeErr(ErrInvalidExpression, node, fmt.Errorf("invald expression output"))
	}

	return ok.(bool), nil
}
