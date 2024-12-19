package ast

import (
	"github.com/bndrmrtn/smarti/internal/lexer"
	"github.com/fatih/color"
)

type Node struct {
	IsReference bool        `json:"is_reference" yaml:"is_reference"`
	Token       lexer.Token `json:"token" yaml:"token"`
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
	Type        NodeType    `json:"type" yaml:"type"`
	Value       string      `json:"value" yaml:"value"`
	Args        []Node      `json:"args,omitempty" yaml:"args,omitempty"`
	Children    []Node      `json:"children,omitempty" yaml:"children,omitempty"`
	Scope       NodeScope   `json:"scope,omitempty" yaml:"scope,omitempty"`

	Info NodeFileInfo `json:"info,omitempty" yaml:"info,omitempty"`
}

type NodeFileInfo struct {
	File string `json:"file" yaml:"file"`
	Pos  int    `json:"pos" yaml:"pos"`
	Line int    `json:"line" yaml:"line"`
}

func (n NodeFileInfo) String() string {
	blue := color.New(color.FgBlue, color.Bold).SprintFunc()

	return blue(n.File) + ":" + blue(n.Line) + ":" + blue(n.Pos)
}

type NodeScope string

const (
	ScopeGlobal NodeScope = "global"
	ScopeLocal  NodeScope = "local"
	ScopeBlock  NodeScope = "block"
	ScopeFunc   NodeScope = "func"
)
