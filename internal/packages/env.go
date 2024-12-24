package packages

import (
	"errors"
	"os"
)

type Env struct{}

func (e Env) Run(fn string, args []*Variable) ([]*FuncReturn, error) {
	switch fn {
	case "get":
		return e.fnGet(args)
	case "set":
		return e.fnSet(args)
	}
	return nil, nil
}

func (Env) fnGet(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 1 {
		return nil, errors.New("get method only allows one argument")
	}

	if args[0].Type != "string" {
		return nil, errors.New("get method only allows string arguments")
	}

	return []*FuncReturn{
		{
			Type:  "string",
			Value: os.Getenv(args[0].Value.(string)),
		},
	}, nil
}

func (Env) fnSet(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 2 {
		return nil, errors.New("set method only allows two argument")
	}

	if args[0].Type != "string" || args[1].Type != "string" {
		return nil, errors.New("set method only allows string arguments")
	}

	os.Setenv(args[0].Value.(string), args[1].Value.(string))

	return nil, nil
}
