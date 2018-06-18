package parser

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/scanner"
	"github.com/patrick-jessen/script/compiler/token"
)

type Parser struct {
	file *file.File

	scanner scanner.Scanner
	tok     token.Token
}

func New(file *file.File) (p *Parser) {
	p = &Parser{
		file: file,
	}
	p.scanner.Init(file)
	p.next()
	return p
}

func (p *Parser) next() {
	if p.tok.ID == token.EOF {
		fmt.Println(p.file.Errors)
		panic("eof")
	}
	p.tok = p.scanner.Scan()
}

func (p *Parser) expect(id token.ID) {
	if p.tok.ID != id {
		p.file.Error(p.tok.Pos, fmt.Sprintf("expected %v", id.String()))
		if len(p.file.Errors) > 10 {
			panic("too many errors")
		}
	} else {
		p.next()
	}
}

func (p *Parser) Debug() {
	fmt.Println(p.file.PosInfo(p.tok.Pos).String())
}

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

func (p *Parser) parseFunctionCall() *ast.FunctionCall {
	ast := &ast.FunctionCall{}

	ast.Identifier = p.parseIdentifier()
	p.expect(token.ParentStart)
	if p.tok.ID != token.ParentEnd {
		ast.Args = p.parseFunctionCallArgs()
	}
	p.expect(token.ParentEnd)

	return ast
}

func (p *Parser) parseStatement() (n ast.Node) {
	switch p.tok.ID {
	case token.Var:
		n = p.parseVariableDecl()
	case token.Identifier:
		n = p.parseFunctionCall()
	default:
		p.file.Error(p.tok.Pos, "expected statement")
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
		stmt := p.parseStatement()
		ast.Statements = append(ast.Statements, stmt)
	}
	p.expect(token.CurlEnd)
	return ast
}

func (p *Parser) parseExpression() ast.Expression {
	// Shunting-yard algorithm
	var val ast.Expression
	var out []interface{}
	var ops []token.ID

loop:
	for {
		switch p.tok.ID {
		case token.String:
			val = p.parseString()
		case token.Identifier:
			val = p.parseIdentifier()
		case token.Integer:
			val = p.parseInteger()
		case token.Float:
			val = p.parseFloat()
		case token.ParentStart:
			p.expect(token.ParentStart)
			val = p.parseExpression()
			p.expect(token.ParentEnd)
		default:
			p.file.Error(p.tok.Pos, "expected expression")
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
			ops = append(ops, p.tok.ID)
			p.expect(p.tok.ID)

		case token.Asterisk:
			fallthrough
		case token.Slash:
			for i := len(ops) - 1; i >= 0; i-- {
				if ops[i] == token.Asterisk || ops[i] == token.Slash {
					out = append(out, ops[i])
					ops = ops[:i]
				}
			}
			ops = append(ops, p.tok.ID)
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
		case token.ID:
			lhs := valStack[len(valStack)-2]
			rhs := valStack[len(valStack)-1]
			valStack = valStack[:len(valStack)-2]

			switch v.(token.ID) {
			case token.Plus:
				expr = &ast.Add{
					LHS: lhs,
					RHS: rhs,
				}
			case token.Minus:
				expr = &ast.Subtract{
					LHS: lhs,
					RHS: rhs,
				}
			case token.Asterisk:
				expr = &ast.Multiply{
					LHS: lhs,
					RHS: rhs,
				}
			case token.Slash:
				expr = &ast.Divide{
					LHS: lhs,
					RHS: rhs,
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
	ast := &ast.Identifier{Token: p.tok}
	p.expect(token.Identifier)
	return ast
}

func (p *Parser) parseFunctionDecl() *ast.FunctionDecl {
	ast := &ast.FunctionDecl{}

	p.expect(token.Func)
	ast.Identifier = p.parseIdentifier()
	p.expect(token.ParentStart)
	p.expect(token.ParentEnd)
	ast.Block = p.parseBlock()

	return ast
}
func (p *Parser) parseVariableDecl() *ast.VariableDecl {
	ast := &ast.VariableDecl{}

	p.expect(token.Var)
	ast.Identifier = p.parseIdentifier()
	p.expect(token.Equal)
	ast.Value = p.parseExpression()

	return ast
}

func (p *Parser) parseDeclaration() (n ast.Node) {
	switch p.tok.ID {
	case token.Func:
		n = p.parseFunctionDecl()
	case token.Var:
		n = p.parseVariableDecl()
	default:
		p.file.Error(p.tok.Pos, "expected declaration")
		p.next()
		return nil
	}
	p.expect(token.NewLine)
	return
}

func (p *Parser) parseFile() *ast.File {
	ast := &ast.File{}

	decl := p.parseDeclaration()
	ast.Declarations = append(ast.Declarations, decl)

	return ast
}

func Run(f *file.File) ast.Node {
	return New(f).parseFile()
}
