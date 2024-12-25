package packages

import (
	"errors"
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
	}
	return nil, nil
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
