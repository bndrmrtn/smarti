package packages

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type IO struct{}

func (i IO) Run(fn string, args []*Variable) ([]*FuncReturn, error) {
	switch fn {
	case "read":
		return i.fnRead(args)
	case "readfile":
		return i.fnReadFile(args)
	case "write":
		return i.fnWrite(args)
	case "writeln":
		return i.fnWrite(args, true)
	case "writef":
		return i.fnWritef(args)
	}
	return nil, fmt.Errorf("function io.%s does not exists", fn)
}

func (IO) Access(variable string) (*Variable, error) {
	return nil, errors.New("io package does not have any variables")
}

func (IO) fnRead(args []*Variable) ([]*FuncReturn, error) {
	var text string
	if len(args) == 0 {
		fmt.Scan(&text)
	} else {
		if args[0].Type == VarString || args[0].Type == VarSingleString {
			fmt.Print(args[0].Value)
			fmt.Scan(&text)
		} else {
			return nil, fmt.Errorf("read expects first argument to be a string")
		}
	}
	return []*FuncReturn{
		{
			Value: text,
			Type:  VarString,
		},
	}, nil
}

func (b IO) fnWrite(args []*Variable, nl ...bool) ([]*FuncReturn, error) {
	values := make([]interface{}, len(args))
	for i, arg := range args {
		values[i] = arg.Value
	}
	if len(nl) > 0 && nl[0] {
		fmt.Println(values...)
	} else {
		fmt.Print(values...)
	}
	return nil, nil
}

func (b IO) fnWritef(args []*Variable) ([]*FuncReturn, error) {
	var format string
	values := make([]interface{}, len(args)-1)
	for i, arg := range args {
		if i == 0 {
			if arg.Type == VarString || arg.Type == VarSingleString {
				format = arg.Value.(string)
			} else {
				return nil, fmt.Errorf("writef expects first argument to be a string")
			}
			continue
		}
		values[i-1] = arg.Value
	}
	fmt.Printf(format, values...)
	return nil, nil
}

func (IO) fnReadFile(args []*Variable) ([]*FuncReturn, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("readfile expects exactly one argument")
	}
	if args[0].Type != VarString {
		return nil, fmt.Errorf("readfile expects string as argument")
	}

	file, err := os.Open(args[0].Value.(string))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return []*FuncReturn{
		{
			Value: string(content),
			Type:  VarString,
		},
	}, nil
}
