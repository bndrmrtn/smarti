package packages

import (
	"errors"
	"fmt"
	"html"
)

type HttpSec struct{}

func (h HttpSec) Run(fn string, args []*Variable) ([]*FuncReturn, error) {
	switch fn {
	case "escapeHTML":
		return h.fnEscapeHTML(args)
	}
	return nil, fmt.Errorf("function httpsec.%s does not exists", fn)
}

func (HttpSec) Access(variable string) (*Variable, error) {
	return nil, errors.New("httpsec package does not have any variables")
}

func (HttpSec) fnEscapeHTML(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 1 {
		return nil, errors.New("escapeHTML requires one argument")
	}

	if args[0].Type != VarString && args[0].Type != VarSingleString {
		return nil, errors.New("escapeHTML requires string arguments")
	}

	escaped := html.EscapeString(args[0].Value.(string))
	return []*FuncReturn{
		{
			Value: escaped,
			Type:  VarString,
		},
	}, nil
}
