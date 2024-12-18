package runtime

import (
	"fmt"
)

type runtimeBuiltin struct{}

var builtin runtimeBuiltin

func (b runtimeBuiltin) runFn(fn string, args []variable) ([]interface{}, error) {
	switch fn {
	case "write":
		return b.runFnWrite(args)
	case "writeln":
		return b.runFnWrite(args, true)
	case "type":
		return builtin.runFnType(args)
	case "writeType":
		return b.runFnWriteType(args)
	}
	return nil, fmt.Errorf("func %s does not exists", fn)
}

func (b runtimeBuiltin) runFnWrite(args []variable, nl ...bool) ([]interface{}, error) {
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

func (b runtimeBuiltin) runFnType(args []variable) ([]interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("type function expects exactly one argument")
	}
	return []interface{}{args[0].Type}, nil
}

func (b runtimeBuiltin) runFnWriteType(args []variable) ([]interface{}, error) {
	var types []interface{}
	for _, arg := range args {
		types = append(types, string(arg.Type))
	}
	fmt.Print(types...)
	return nil, nil
}
