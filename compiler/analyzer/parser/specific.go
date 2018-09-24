package parser

import (
	"github.com/patrick-jessen/script/compiler/analyzer/ast"
	"github.com/patrick-jessen/script/utils/token"
)

func (p *Parser) parseFunctionCallArgs() *ast.FunctionCallArgs {
	ast := &ast.FunctionCallArgs{}

	for {
		expr := p.parseExpression()
		ast.Args = append(ast.Args, expr)

		if p.tok.ID != token.Comma {
			break
		}
		p.expect(token.Comma)
	}

	return ast
}

func (p *Parser) parseFunctionCall(ident *ast.Identifier) *ast.FunctionCall {
	ast := &ast.FunctionCall{}

	ast.Identifier = ident
	p.expect(token.ParentStart)
	if p.tok.ID != token.ParentEnd {
		ast.Args = p.parseFunctionCallArgs()
	}
	ast.LastParentPos = p.tok.Pos
	p.expect(token.ParentEnd)
	p.Resolve(ast.Identifier)
	return ast
}

func (p *Parser) parseVariableAssign(ident *ast.Identifier) *ast.VariableAssign {
	ast := &ast.VariableAssign{}

	ast.Identifier = ident
	ast.EqPos = p.tok.Pos
	p.expect(token.Equal)
	ast.Value = p.parseExpression()

	p.Resolve(ast.Identifier)
	return ast
}

func (p *Parser) parseStatement() (n ast.Node) {
	switch p.tok.ID {
	case token.Var:
		n = p.parseVariableDecl()
	case token.Identifier:
		ident := p.parseIdentifier()
		if p.tok.ID == token.ParentStart {
			n = p.parseFunctionCall(ident)
		} else {
			n = p.parseVariableAssign(ident)
		}
	default:
		p.tok.Pos.MarkError("expected statement")
		p.next()
		return nil
	}
	p.expect(token.NewLine)
	return
}

func (p *Parser) parseBlock() *ast.Block {
	ast := &ast.Block{}

	p.expect(token.CurlStart)
	p.expect(token.NewLine)
	for p.tok.ID != token.CurlEnd {
		if p.tok.ID == token.NewLine {
			p.next()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			ast.Statements = append(ast.Statements, stmt)
		}
	}
	p.expect(token.CurlEnd)
	return ast
}

func (p *Parser) parseExpression() ast.Expression {
	// Shunting-yard algorithm
	var val ast.Expression
	var out []interface{}
	var ops []token.Token

loop:
	for {
		switch p.tok.ID {
		case token.String:
			val = p.parseString()
		case token.Identifier:
			ident := p.parseIdentifier()
			ref := &ast.VariableRef{Identifier: ident}
			val = ref
			p.Resolve(ref.Identifier)
		case token.Integer:
			val = p.parseInteger()
		case token.Float:
			val = p.parseFloat()
		case token.ParentStart:
			p.expect(token.ParentStart)
			val = p.parseExpression()
			p.expect(token.ParentEnd)
		default:
			p.tok.Pos.MarkError("expected expression")
			p.next()
			return nil
		}
		out = append(out, val)

		switch p.tok.ID {
		case token.Plus:
			fallthrough
		case token.Minus:
			for i := len(ops) - 1; i >= 0; i-- {
				out = append(out, ops[i])
				ops = ops[:i]
			}
			ops = append(ops, p.tok)
			p.expect(p.tok.ID)

		case token.Asterisk:
			fallthrough
		case token.Slash:
			for i := len(ops) - 1; i >= 0; i-- {
				if ops[i].ID == token.Asterisk || ops[i].ID == token.Slash {
					out = append(out, ops[i])
					ops = ops[:i]
				}
			}
			ops = append(ops, p.tok)
			p.expect(p.tok.ID)
		default:
			break loop
		}
	}

	for i := len(ops) - 1; i >= 0; i-- {
		out = append(out, ops[i])
	}

	var valStack []ast.Expression
	var expr ast.Expression

	for _, v := range out {
		switch v.(type) {
		case nil:
			return nil
		case token.Token:
			lhs := valStack[len(valStack)-2]
			rhs := valStack[len(valStack)-1]
			valStack = valStack[:len(valStack)-2]

			tok := v.(token.Token)
			switch tok.ID {
			case token.Plus:
				expr = &ast.Add{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			case token.Minus:
				expr = &ast.Subtract{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			case token.Asterisk:
				expr = &ast.Multiply{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			case token.Slash:
				expr = &ast.Divide{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			}

			valStack = append(valStack, expr)
		default:
			valStack = append(valStack, v.(ast.Expression))
		}
	}

	return valStack[0]
}

func (p *Parser) parseInteger() *ast.Integer {
	ast := &ast.Integer{Token: p.tok}
	p.expect(token.Integer)
	return ast
}

func (p *Parser) parseFloat() *ast.Float {
	ast := &ast.Float{Token: p.tok}
	p.expect(token.Float)
	return ast
}

func (p *Parser) parseString() *ast.String {
	ast := &ast.String{Token: p.tok}
	p.expect(token.String)
	return ast
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	ast := &ast.Identifier{}

	ident := p.tok
	p.expect(token.Identifier)

	if p.tok.ID == token.Dot {
		p.expect(token.Dot)
		ast.Module = ident
		ast.Symbol = p.tok
		p.expect(token.Identifier)
	} else {
		ast.Symbol = ident
	}

	return ast
}

func (p *Parser) parseFunctionDeclArgs() []*ast.Identifier {
	var ret []*ast.Identifier

	for {
		ident := p.parseIdentifier()
		ident.Obj = &ast.Object{}
		ident.Typ = ast.Type{
			IsResolved: true,
			Return:     p.tok.Value,
		}
		p.expect(token.Identifier)

		ret = append(ret, ident)

		if p.tok.ID != token.Comma {
			break
		}
		p.expect(token.Comma)
	}
	return ret
}

func (p *Parser) parseType() ast.Type {
	typ := ast.Type{
		IsResolved: true,
	}

	typ.Return = p.tok.Value
	p.expect(token.Identifier)
	return typ
}

func (p *Parser) parseFunctionDecl() *ast.FunctionDecl {
	obj := &ast.Object{}
	ast := &ast.FunctionDecl{}

	p.expect(token.Func)
	ast.Identifier = p.parseIdentifier()
	ast.Identifier.Obj = obj
	p.expect(token.ParentStart)
	if p.tok.ID != token.ParentEnd {
		ast.Args = p.parseFunctionDeclArgs()
	}
	p.expect(token.ParentEnd)

	p.pushScope()
	for _, a := range ast.Args {
		p.Declare(a)
	}
	ast.Block = p.parseBlock()
	p.popScope()

	ast.Init()

	p.Declare(ast)
	return ast
}
func (p *Parser) parseVariableDecl() *ast.VariableDecl {
	obj := &ast.Object{}
	ast := &ast.VariableDecl{}

	p.expect(token.Var)
	ast.Identifier = p.parseIdentifier()
	ast.Identifier.Obj = obj
	if p.tok.ID != token.Equal {
		ast.Identifier.Typ = p.parseType()
	}
	p.expect(token.Equal)
	ast.Value = p.parseExpression()

	p.Declare(ast)
	return ast
}

func (p *Parser) parseDeclaration() (n ast.Node) {
	switch p.tok.ID {
	case token.Func:
		n = p.parseFunctionDecl()
	case token.Var:
		n = p.parseVariableDecl()
	default:
		p.tok.Pos.MarkError("expected declaration")
		p.next()
		return nil
	}
	p.expect(token.NewLine)
	return
}

func (p *Parser) parseImport() {
	var alias token.Token
	var module token.Token

	p.expect(token.Import)
	if p.tok.ID == token.Identifier {
		alias = p.tok
		p.expect(token.Identifier)
	}
	module = p.tok
	p.expect(token.String)
	p.expect(token.NewLine)
	p.importModule(Import{alias, module})
}

func (p *Parser) parseFile() *ast.File {
	ast := &ast.File{}

	for p.tok.ID == token.Import {
		p.parseImport()
		if p.tok.ID == token.NewLine {
			p.expect(token.NewLine)
		}
	}

	for p.tok.ID != token.EOF {
		if p.tok.ID == token.NewLine {
			p.next()
			continue
		}
		decl := p.parseDeclaration()
		if decl != nil {
			ast.Declarations = append(ast.Declarations, decl)
		}
	}

	return ast
}
