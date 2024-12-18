package ast

type NodeType string

const (
	VarNil          NodeType = "nil"
	VarString       NodeType = "string"
	VarSingleString NodeType = "string_single"
	VarNumber       NodeType = "number"
	VarFloat        NodeType = "float"
	VarBool         NodeType = "bool"
	VarTemplate     NodeType = "template"
	VarVariable     NodeType = "variable"

	VarUnknown NodeType = "#unknown#"

	FuncCall NodeType = "func_call"
)
