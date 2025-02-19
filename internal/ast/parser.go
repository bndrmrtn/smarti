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
			pkgToken := p.tokens[inx]
			pkg := pkgToken.Value
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
				Info:  getInfo(pkgToken),
			})
		case lexer.Namespace:
			name := p.tokens[inx].Value
			info := getInfo(p.tokens[inx])
			inx++
			if inx < tokenLen && p.tokens[inx].Type != lexer.SemiColon {
				return errors.New("syntax error: missing semicolon")
			}
			p.Nodes = append(p.Nodes, Node{
				Token: lexer.Namespace,
				Type:  Namespace,
				Name:  name,
				Info:  info,
			})
		case lexer.Func:
			name, args := getFuncCall(p.tokens[inx])
			info := getInfo(p.tokens[inx])
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
				Info:     info,
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
		case lexer.If:
			var (
				conditionTokens []lexer.LexerToken
				bodyTokens      []lexer.LexerToken
			)

			for inx < tokenLen && p.tokens[inx].Type != lexer.CurlyBraceStart {
				conditionTokens = append(conditionTokens, p.tokens[inx])
				inx++
			}

			if inx >= tokenLen || p.tokens[inx].Type != lexer.CurlyBraceStart {
				return errors.New("syntax error: missing opening curly brace for if statement")
			}

			for inx >= tokenLen || p.tokens[inx].Type != lexer.CurlyBraceStart {
				inx++
			}

			depth := 0

			if inx >= tokenLen || p.tokens[inx].Type != lexer.CurlyBraceStart {
				return errors.New("syntax error: missing opening curly brace")
			}

			depth++
			inx++

			for inx < tokenLen {
				token := p.tokens[inx]

				if token.Type == lexer.CurlyBraceStart {
					depth++
				} else if token.Type == lexer.CurlyBraceEnd {
					depth--
					if depth == 0 {
						inx++
						break
					}
				}

				if depth > 0 {
					bodyTokens = append(bodyTokens, token)
				}

				inx++
			}

			if depth != 0 {
				return errors.New("syntax error: unbalanced curly braces")
			}

			var condition Node
			bindValue(conditionTokens, &condition)

			bodyParser := NewParser(bodyTokens)
			if err := bodyParser.Parse(); err != nil {
				return err
			}

			p.Nodes = append(p.Nodes, Node{
				Token:    lexer.If,
				Type:     IfStatement,
				Args:     []Node{condition},
				Children: append([]Node{condition}, bodyParser.Nodes...),
			})
		case lexer.For:
			var (
				initTokens      []lexer.LexerToken
				conditionTokens []lexer.LexerToken
				postTokens      []lexer.LexerToken
				bodyTokens      []lexer.LexerToken
			)

			// Első szakasz: init
			for inx < tokenLen && p.tokens[inx].Type != lexer.SemiColon {
				initTokens = append(initTokens, p.tokens[inx])
				inx++
			}

			if inx >= tokenLen || p.tokens[inx].Type != lexer.SemiColon {
				return errors.New("syntax error: missing semicolon after init statement in for loop")
			}
			inx++ // Átlépünk a ';' után

			// Második szakasz: condition
			for inx < tokenLen && p.tokens[inx].Type != lexer.SemiColon {
				conditionTokens = append(conditionTokens, p.tokens[inx])
				inx++
			}

			if inx >= tokenLen || p.tokens[inx].Type != lexer.SemiColon {
				return errors.New("syntax error: missing semicolon after condition statement in for loop")
			}
			inx++ // Átlépünk a ';' után

			// Harmadik szakasz: post
			for inx < tokenLen && p.tokens[inx].Type != lexer.CurlyBraceStart {
				postTokens = append(postTokens, p.tokens[inx])
				inx++
			}

			if inx >= tokenLen || p.tokens[inx].Type != lexer.CurlyBraceStart {
				return errors.New("syntax error: missing opening curly brace for for loop body")
			}
			inx++ // Átlépünk a '{' után

			// Negyedik szakasz: body
			depth := 1
			for inx < tokenLen && depth > 0 {
				if p.tokens[inx].Type == lexer.CurlyBraceStart {
					depth++
				} else if p.tokens[inx].Type == lexer.CurlyBraceEnd {
					depth--
					if depth == 0 {
						break
					}
				}

				if depth > 0 {
					bodyTokens = append(bodyTokens, p.tokens[inx])
				}
				inx++
			}

			if depth != 0 {
				return errors.New("syntax error: unbalanced curly braces in for loop body")
			}

			// Parsoljuk az init, condition és post részeket
			var initNode, conditionNode, postNode Node
			bindValue(initTokens, &initNode)
			bindValue(conditionTokens, &conditionNode)
			bindValue(postTokens, &postNode)

			// Parsoljuk a body-t
			bodyParser := NewParser(bodyTokens)
			if err := bodyParser.Parse(); err != nil {
				return err
			}

			// Létrehozzuk a for loop node-ot
			p.Nodes = append(p.Nodes, Node{
				Token:    lexer.For,
				Type:     ForLoop,
				Args:     []Node{initNode, conditionNode, postNode},
				Children: bodyParser.Nodes,
			})

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
