package packages

type Variable struct {
	Type  VarType
	Value interface{}
}

type FuncReturn struct {
	Type  VarType
	Value interface{}
}

type Package interface {
	Run(fn string, args []*Variable) ([]*FuncReturn, error)
}
