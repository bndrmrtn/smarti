package packages

import (
	"errors"
	"fmt"
	"net/http"
)

type Request struct {
	r *http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		r: r,
	}
}

func (r *Request) Run(fn string, args []*Variable) ([]*FuncReturn, error) {
	switch fn {
	case "method":
		return r.fnMethod(args)
	case "query":
		return r.fnQuery(args)
	}
	return nil, fmt.Errorf("function request.%s does not exists", fn)
}

func (*Request) Access(variable string) (*Variable, error) {
	return nil, errors.New("request package does not have any variables")
}

func (r *Request) fnMethod(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 0 {
		return nil, errors.New("method does not accept any argument")
	}

	return []*FuncReturn{
		{
			Value: r.r.Method,
			Type:  VarString,
		},
	}, nil
}

func (r *Request) fnQuery(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 1 {
		return nil, errors.New("query requires exatly one argument")
	}

	if args[0].Type != VarString && args[0].Type != VarSingleString {
		return nil, errors.New("query requires string arguments")
	}

	return []*FuncReturn{
		{
			Value: r.r.URL.Query().Get(args[0].Value.(string)),
			Type:  VarString,
		},
	}, nil
}
