package packages

type VarType string

const (
	VarNil          VarType = "nil"
	VarString       VarType = "string"
	VarSingleString VarType = "string_single"
	VarNumber       VarType = "number"
	VarFloat        VarType = "float"
	VarBool         VarType = "bool"
	VarTemplate     VarType = "template"
	VarVariable     VarType = "variable"

	VarUnknown VarType = "#unknown#"

	FuncCall VarType = "func_call"

	VarExpression VarType = "expression"
	VarOperator   VarType = "operator"

	UsePackage VarType = "use_package"
)
