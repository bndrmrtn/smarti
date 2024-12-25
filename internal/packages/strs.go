package packages

import (
	"errors"
	"fmt"
	"strings"
)

type Strs struct{}

func (s Strs) Run(fn string, args []*Variable) ([]*FuncReturn, error) {
	switch fn {
	case "length":
		return s.fnLength(args)
	case "trim":
		return s.fnTrim(args)
	case "concat":
		return s.fnConcat(args)
	}
	return nil, fmt.Errorf("function strs.%s does not exists", fn)
}

func (Strs) Access(variable string) (*Variable, error) {
	return nil, errors.New("strs package does not have any variables")
}

func (Strs) fnLength(args []*Variable) ([]*FuncReturn, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("length expects at least one argument")
	}
	arg := args[0]
	if arg.Type != VarString && arg.Type != VarSingleString {
		return nil, fmt.Errorf("length expects a string argument")
	}

	l := len(arg.Value.(string))
	return []*FuncReturn{
		{
			Value: l,
			Type:  VarNumber,
		},
	}, nil
}

func (Strs) fnTrim(args []*Variable) ([]*FuncReturn, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("trim expects at least one argument")
	}
	arg := args[0]
	if arg.Type != VarString && arg.Type != VarSingleString {
		return nil, fmt.Errorf("trim expects a string argument")
	}

	return []*FuncReturn{
		{
			Value: strings.TrimSpace(arg.Value.(string)),
			Type:  VarString,
		},
	}, nil
}

func (Strs) fnConcat(args []*Variable) ([]*FuncReturn, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("concat expects at least one argument")
	}

	var sb strings.Builder

	for _, arg := range args {
		if arg.Type != VarString && arg.Type != VarSingleString {
			return nil, fmt.Errorf("concat expects only string argument")
		}

		sb.WriteString(arg.Value.(string))
	}

	return []*FuncReturn{
		{
			Value: sb.String(),
			Type:  VarString,
		},
	}, nil
}
