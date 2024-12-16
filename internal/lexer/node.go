package lexer

type Node struct {
	IsReference bool
	Token       Token
	Name        string
	Type        string
	Value       interface{}
	Args        []Node
	Children    []Node

	Info NodeFileInfo
}

type NodeFileInfo struct {
	File string
	Pos  int
	Line int
}
