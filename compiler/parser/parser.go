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
	p.tok = p.scanner.Scan()
}

func (p *Parser) expect(id token.ID) {
	if p.tok.ID != id {
		p.file.Error(p.tok.Pos, fmt.Sprintf("expected %v", id.String()))
		if len(p.file.Errors) > 10 {
			panic("too many errors")
		}
	}
	p.next()
}

func (p *Parser) Debug() {
	fmt.Println(p.file.PosInfo(p.tok.Pos).String())
}

func (p *Parser) parseFunctionCallArgs() *ast.FunctionCallArgs {
	ast := &ast.FunctionCallArgs{}

	switch p.tok.ID {
	case token.Identifier:
		ident := p.parseIdentifier()
		ast.Args = append(ast.Args, ident)
	default:
		p.file.Error(p.tok.Pos, "expected argument")
		p.next()
		return nil
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
	var exp ast.Expression

	switch p.tok.ID {
	case token.String:
		exp = p.parseString()
	default:
		p.file.Error(p.tok.Pos, "expected declaration")
		p.next()
		return nil
	}

	return exp
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
