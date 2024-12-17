package ast

import "github.com/smlgh/smarti/internal/lexer"

type Node struct {
	IsReference bool        `json:"is_reference"`
	Token       lexer.Token `json:"token"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
	Args        []Node      `json:"args"`
	Children    []Node      `json:"children"`

	Info NodeFileInfo `json:"info"`
}

type NodeFileInfo struct {
	File string `json:"file"`
	Pos  int    `json:"pos"`
	Line int    `json:"line"`
}
