package runtime

import (
	"fmt"
	"io"
	"os"

	"github.com/smlgh/smarti/internal/ast"
)

type IO struct{}

func (i IO) Run(fn string, args []variable) ([]funcReturn, error) {
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
	return nil, nil
}

func (IO) fnRead(args []variable) ([]funcReturn, error) {
	var text string
	if len(args) == 0 {
		fmt.Scan(&text)
	} else {
		if args[0].Type == ast.VarString || args[0].Type == ast.VarSingleString {
			fmt.Print(args[0].Value)
			fmt.Scan(&text)
		} else {
			return nil, fmt.Errorf("read expects first argument to be a string")
		}
	}
	return []funcReturn{
		{
			Value: text,
			Type:  ast.VarString,
		},
	}, nil
}

func (b IO) fnWrite(args []variable, nl ...bool) ([]funcReturn, error) {
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

func (b IO) fnWritef(args []variable) ([]funcReturn, error) {
	var format string
	values := make([]interface{}, len(args)-1)
	for i, arg := range args {
		if i == 0 {
			if arg.Type == ast.VarString || arg.Type == ast.VarSingleString {
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

func (IO) fnReadFile(args []variable) ([]funcReturn, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("readfile expects exactly one argument")
	}
	if args[0].Type != ast.VarString {
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

	return []funcReturn{
		{
			Value: string(content),
			Type:  ast.VarString,
		},
	}, nil
}
