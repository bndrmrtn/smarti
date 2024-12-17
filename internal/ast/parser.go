package ast

import (
	"fmt"

	"github.com/smlgh/smarti/internal/lexer"
)

type Parser struct {
	tokens []lexer.LexerToken

	Nodes []Node
}

func NewParser(tokens []lexer.LexerToken) *Parser {
	return &Parser{
		tokens: tokens,
		Nodes:  make([]Node, 0),
	}
}

func (p *Parser) Parse() error {
	return p.parse()
}

func (p *Parser) parse() error {
	var (
		inx      = 0
		tokenLen = len(p.tokens)
	)

	for inx+1 < tokenLen {
		token := p.tokens[inx]
		inx++

		switch token.Type {
		case lexer.Let, lexer.Const:
			info := getInfo(token)
			value, typ, ref := getType(p.tokens[inx])
			if err := p.canAssign(token.Value, true); err != nil {
				return NewErrWithPos(info, err.Error())
			}
			n := Node{
				IsReference: ref,
				Token:       token.Type,
				Name:        token.Value,
				Value:       value,
				Type:        typ,
				Info:        info,
			}
			p.Nodes = append(p.Nodes, n)
			inx++
			continue
		case lexer.Assign:
			info := getInfo(token)
			value, typ, ref := getType(token)
			variable := p.tokens[inx-2]
			if err := p.canAssign(variable.Value, false); err != nil {
				return NewErrWithPos(info, err.Error())
			}
			n := Node{
				IsReference: ref,
				Token:       token.Type,
				Name:        variable.Value,
				Value:       value,
				Type:        typ,
				Info:        info,
			}
			p.Nodes = append(p.Nodes, n)
			continue
		case lexer.FuncCall:
			info := getInfo(token)
			fmt.Println(info)
		}
	}

	return nil
}

func (p *Parser) canAssign(name string, create bool) error {
	for _, n := range p.Nodes {
		if !create && n.Name == name && n.Token == lexer.Const {
			return ErrorCannotReAssignConst
		}

		if create && n.Name == name {
			return ErrorCannotReDeclareVar
		}
	}

	return nil
}
