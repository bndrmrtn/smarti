package packages

import (
	"errors"
	"fmt"
	"strconv"
)

type Numbers struct{}

func (r Numbers) Run(fn string, args []*Variable) ([]*FuncReturn, error) {
	switch fn {
	case "from":
		return r.fnFrom(args)
	}
	return nil, fmt.Errorf("function numbers.%s does not exists", fn)
}

func (Numbers) Access(variable string) (*Variable, error) {
	return nil, errors.New("numbers package does not have any variables")
}

func (Numbers) fnFrom(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 1 {
		return nil, errors.New("from expects one argument")
	}

	switch args[0].Type {
	case VarString:
		val, err := strconv.Atoi(args[0].Value.(string))
		if err != nil {
			return nil, err
		}

		return []*FuncReturn{
			{
				Value: val,
				Type:  VarNumber,
			},
		}, nil
	case VarNumber:
		return []*FuncReturn{
			{
				Value: args[0].Value,
				Type:  VarNumber,
			},
		}, nil
	case VarFloat:
		return []*FuncReturn{
			{
				Value: int(args[0].Value.(float64)),
				Type:  VarNumber,
			},
		}, nil
	}

	return nil, errors.New("from expects a string, int or float")
}
