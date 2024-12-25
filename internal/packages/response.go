package packages

import (
	"errors"
	"fmt"
	"net/http"
)

type Response struct {
	rw http.ResponseWriter
}

func NewResponse(rw http.ResponseWriter) *Response {
	return &Response{
		rw: rw,
	}
}

func (r *Response) Run(fn string, args []*Variable) ([]*FuncReturn, error) {
	switch fn {
	case "write":
		return r.fnWrite(args)
	case "header":
		return r.fnHeader(args)
	case "status":
		return r.fnStatus(args)
	}
	return nil, fmt.Errorf("function response.%s does not exists", fn)
}

func (*Response) Access(variable string) (*Variable, error) {
	return nil, errors.New("response package does not have any variables")
}

func (r *Response) fnWrite(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 1 {
		return nil, errors.New("write method only accepts one argument")
	}

	if args[0].Type == VarString || args[0].Type == VarSingleString {
		r.rw.Write([]byte(args[0].Value.(string)))
		return nil, nil
	}

	return nil, errors.New("write method only accepts string argument")
}

func (r *Response) fnHeader(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 2 {
		return nil, errors.New("header method accepts 2 arguments, key and value")
	}

	if args[0].Type == VarString || args[0].Type == VarSingleString &&
		args[1].Type == VarString || args[1].Type == VarSingleString {
		r.rw.Header().Add(args[0].Value.(string), args[1].Value.(string))
		return nil, nil
	}

	return nil, errors.New("header method only accepts string arguments")
}

func (r *Response) fnStatus(args []*Variable) ([]*FuncReturn, error) {
	if len(args) != 1 {
		return nil, errors.New("status method accepts 1 arguments")
	}

	if args[0].Type == VarNumber {
		r.rw.WriteHeader(args[0].Value.(int))
		return nil, nil
	}

	return nil, errors.New("status method only accepts number arguments")
}
