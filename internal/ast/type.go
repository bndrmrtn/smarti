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

	FuncCall   NodeType = "func_call"
	FuncDecl   NodeType = "func_decl"
	FuncReturn NodeType = "func_return"

	VarExpression NodeType = "expression"
	VarOperator   NodeType = "operator"

	UsePackage NodeType = "use_package"
	Namespace  NodeType = "namespace"

	ForLoop     NodeType = "for_loop"
	IfStatement NodeType = "if_statement"
)
