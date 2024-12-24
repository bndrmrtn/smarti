package ast

import (
	"errors"

	"github.com/bndrmrtn/smarti/internal/lexer"
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

	for inx < tokenLen {
		token := p.tokens[inx]
		inx++

		switch token.Type {
		case lexer.Let, lexer.Const:
			name := p.tokens[inx].Value
			value := []lexer.LexerToken{}
			inx++
			if inx < tokenLen && p.tokens[inx].Type == lexer.Assign {
				inx++
				for inx < tokenLen && p.tokens[inx].Type != lexer.SemiColon {
					value = append(value, p.tokens[inx])
					inx++
				}
			} else {
				value = append(value, lexer.LexerToken{
					Type:  lexer.Nil,
					Value: "nil",
				})
			}

			if err := p.canAssign(name, true); err != nil {
				return err
			}

			n := Node{
				Token: token.Type,
				Name:  name,
			}

			bindValue(value, &n)

			p.Nodes = append(p.Nodes, n)
			continue
		case lexer.Assign:
			if inx-2 < 0 {
				return errors.New("syntax error: missing variable name")
			}
			name := p.tokens[inx-2].Value
			value := []lexer.LexerToken{}

			for inx < tokenLen && p.tokens[inx].Type != lexer.SemiColon {
				value = append(value, p.tokens[inx])
				inx++
			}

			if err := p.canAssign(name, false); err != nil {
				return err
			}

			n := Node{
				Token: lexer.Assign,
				Name:  name,
			}

			bindValue(value, &n)

			p.Nodes = append(p.Nodes, n)
			continue
		case lexer.Use:
			pkg := p.tokens[inx].Value
			as := pkg
			if inx+3 < tokenLen && p.tokens[inx+1].Value == "as" {
				inx += 2
				as = p.tokens[inx].Value
				if p.tokens[inx+1].Type != lexer.SemiColon {
					return errors.New("syntax error: missing semicolon")
				}
			} else if inx+1 < tokenLen && p.tokens[inx+1].Type != lexer.SemiColon {
				return errors.New("syntax error: missing semicolon")
			}
			p.Nodes = append(p.Nodes, Node{
				Token: lexer.Use,
				Type:  UsePackage,
				Name:  pkg,
				Value: as,
			})
		case lexer.Namespace:
			name := p.tokens[inx].Value
			inx++
			if inx < tokenLen && p.tokens[inx].Type != lexer.SemiColon {
				return errors.New("syntax error: missing semicolon")
			}
			p.Nodes = append(p.Nodes, Node{
				Token: lexer.Namespace,
				Type:  Namespace,
				Name:  name,
			})
		case lexer.Func:
			name, args := getFuncCall(p.tokens[inx])
			inx++

			body := []lexer.LexerToken{}
			for inx < tokenLen && p.tokens[inx].Type != lexer.CurlyBraceEnd {
				body = append(body, p.tokens[inx])
				inx++
			}

			psr := NewParser(body)
			if err := psr.Parse(); err != nil {
				return err
			}

			for j := range psr.Nodes {
				psr.Nodes[j].Scope = ScopeFunc
			}

			p.Nodes = append(p.Nodes, Node{
				Token:    lexer.Func,
				Type:     FuncDecl,
				Children: psr.Nodes,
				Name:     name,
				Args:     args,
			})
		case lexer.FuncCall:
			name, args := getFuncCall(token)
			p.Nodes = append(p.Nodes, Node{
				Token: lexer.FuncCall,
				Type:  FuncCall,
				Name:  name,
				Args:  args,
				Info:  getInfo(token),
			})
		case lexer.Return:
			returnsRaw := []lexer.LexerToken{}
			for inx < tokenLen && p.tokens[inx].Type != lexer.SemiColon {
				returnsRaw = append(returnsRaw, p.tokens[inx])
				inx++
			}

			var returns Node
			bindValue(returnsRaw, &returns)

			n := Node{
				Token:    lexer.Return,
				Type:     FuncReturn,
				Children: []Node{returns},
			}
			p.Nodes = append(p.Nodes, n)
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

func bindValue(value []lexer.LexerToken, n *Node) {
	if len(value) == 0 {
		n.Value = "nil"
		n.Type = VarNil
		return
	}

	if len(value) == 1 {
		if value[0].Type == lexer.FuncCall {
			name, args := getFuncCall(value[0])
			n.Children = append(n.Children, Node{
				Type: FuncCall,
				Name: name,
				Args: args,
			})
			n.Type = VarExpression
			return
		}

		val, typ, ref := getType(value[0])
		n.Value = val
		n.Type = typ
		n.IsReference = ref
		return
	}

	n.Type = VarExpression
	inx := 0
	for inx < len(value) {
		v := value[inx]
		inx++

		if v.Type == lexer.FuncCall {
			name, args := getFuncCall(v)
			n.Children = append(n.Children, Node{
				Type: FuncCall,
				Name: name,
				Args: args,
				Info: getInfo(v),
			})
			continue
		}

		val, typ, ref := getType(v)
		n.Children = append(n.Children, Node{
			Value:       val,
			Type:        typ,
			IsReference: ref,
			Info:        getInfo(v),
		})
	}
}
