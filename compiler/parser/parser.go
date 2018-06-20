package parser

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/scanner"
	"github.com/patrick-jessen/script/compiler/token"
)

func (p *Parser) Symbols() map[string]ast.Declarable {
	return p.rootScope.symbols
}

type Link struct {
	Decl ast.Declarable
	Ref  ast.Resolvable
}

type Scope struct {
	parent  *Scope
	symbols map[string]ast.Declarable
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]ast.Declarable),
	}
}

func (p *Parser) Declare(d ast.Declarable) {
	name := d.Name()
	sym, ok := p.curScope.symbols[name]
	if ok {
		p.file.Error(d.Pos(), fmt.Sprintf(
			"redeclaration of symbol '%v'. First declared here: (%v)",
			name, p.file.PosInfo(sym.Pos()).Link(),
		))
		return
	}
	p.curScope.symbols[name] = d
}
func (p *Parser) Resolve(r ast.Resolvable) {
	name := r.Name()
	scope := p.curScope

	for scope != nil {
		sym, ok := scope.symbols[name]
		if ok {
			r.SetType(sym.Type())
			return
		}
		scope = scope.parent
	}
	p.Unresolved = append(p.Unresolved, Link{Ref: r})
}

type Parser struct {
	file *file.File

	scanner scanner.Scanner
	tok     token.Token

	rootScope *Scope
	curScope  *Scope

	Unresolved []Link
}

func New(file *file.File) (p *Parser) {
	p = &Parser{
		file:      file,
		rootScope: NewScope(nil),
	}
	return p
}

func (p *Parser) pushScope() {
	p.curScope = NewScope(p.curScope)
}
func (p *Parser) popScope() {
	p.curScope = p.curScope.parent
}

func (p *Parser) next() {
	if p.tok.ID == token.EOF {
		panic("EOF reached")
	}
	p.tok = p.scanner.Scan()
}

func (p *Parser) expect(id token.ID) {
	if p.tok.ID != id {
		p.file.Error(p.tok.Pos, fmt.Sprintf("expected %v", id.String()))
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

func (p *Parser) parseFunctionCall(ident *ast.Identifier) *ast.FunctionCall {
	ast := &ast.FunctionCall{}
	// numArgs := 0

	ast.Identifier = ident
	p.expect(token.ParentStart)
	if p.tok.ID != token.ParentEnd {
		ast.Args = p.parseFunctionCallArgs()
		// numArgs = len(ast.Args.Args)
	}
	// parentEndPos := p.tok.Pos
	p.expect(token.ParentEnd)

	p.Resolve(ast)

	// if ast.Type() == "" {
	// 	return ast
	// }
	// tsplit := strings.Split(ast.Type(), " -> ")
	// argsstr := tsplit[0]
	// args := strings.Split(argsstr[1:len(argsstr)-1], ",")

	// for i := 0; i < numArgs; i++ {
	// 	a := ast.Args.Args[i]
	// 	if len(args) <= i {
	// 		p.file.Error(a.Pos(), fmt.Sprintf(
	// 			"too many arguments, expected %v got %v", len(args), numArgs,
	// 		))
	// 		break
	// 	}
	// 	if args[i] != a.Type() {
	// 		typ := a.Type()
	// 		if typ == "" {
	// 			typ = "undeclared"
	// 		}
	// 		p.file.Error(a.Pos(), fmt.Sprintf(
	// 			"cannot use %v as %v", typ, args[i],
	// 		))
	// 	}
	// }
	// if len(args) > numArgs {
	// 	p.file.Error(parentEndPos, fmt.Sprintf(
	// 		"too few arguments, expected %v got %v", len(args), numArgs,
	// 	))
	// }

	return ast
}

func (p *Parser) parseVariableAssign(ident *ast.Identifier) *ast.VariableAssign {
	ast := &ast.VariableAssign{}

	ast.Identifier = ident
	p.expect(token.Equal)
	ast.Value = p.parseExpression()

	p.Resolve(ast)
	// if ast.Type() != ast.Value.Type() {
	// 	p.file.Error(ast.Value.Pos(), fmt.Sprintf(
	// 		"cannot assign '%v' (type %v) with type %v",
	// 		ast.Name(), ast.Type(), ast.Value.Type()),
	// 	)
	// }

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
			p.Resolve(ref)
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

			typ := lhs.Type()
			var errMsg string

			switch v.(token.Token).ID {
			case token.Plus:
				expr = &ast.Add{
					LHS: lhs,
					RHS: rhs,
					Typ: typ,
				}
				errMsg = "add"
			case token.Minus:
				expr = &ast.Subtract{
					LHS: lhs,
					RHS: rhs,
					Typ: typ,
				}
				errMsg = "subtract"
			case token.Asterisk:
				expr = &ast.Multiply{
					LHS: lhs,
					RHS: rhs,
					Typ: typ,
				}
				errMsg = "multiply"
			case token.Slash:
				expr = &ast.Divide{
					LHS: lhs,
					RHS: rhs,
					Typ: typ,
				}
				errMsg = "divide"
			}

			if lhs.Type().Return != rhs.Type().Return {
				p.file.Error(v.(token.Token).Pos, fmt.Sprintf(
					"cannot %v types %v and %v", errMsg, lhs.Type(), rhs.Type(),
				))
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

func (p *Parser) parseFunctionDeclArgs() *ast.FunctionDeclArgs {
	ast := &ast.FunctionDeclArgs{}

	for {
		ident := p.parseIdentifier()
		typ := p.parseIdentifier()
		ast.Names = append(ast.Names, ident)
		ast.Types = append(ast.Types, typ)

		if p.tok.ID != token.Comma {
			break
		}
		p.expect(token.Comma)
	}
	return ast
}

func (p *Parser) parseFunctionDecl() *ast.FunctionDecl {
	ast := &ast.FunctionDecl{}

	p.expect(token.Func)
	ast.Identifier = p.parseIdentifier()
	p.expect(token.ParentStart)
	if p.tok.ID != token.ParentEnd {
		ast.Args = p.parseFunctionDeclArgs()
	}
	p.expect(token.ParentEnd)

	p.pushScope()
	ast.Block = p.parseBlock()
	p.popScope()

	p.Declare(ast)
	return ast
}
func (p *Parser) parseVariableDecl() *ast.VariableDecl {
	ast := &ast.VariableDecl{}

	p.expect(token.Var)
	ast.Identifier = p.parseIdentifier()
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
		p.file.Error(p.tok.Pos, "expected declaration")
		p.next()
		return nil
	}
	p.expect(token.NewLine)
	return
}

func (p *Parser) parseFile() *ast.File {
	ast := &ast.File{}

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

func (p *Parser) Run() ast.Node {
	p.tok = token.Token{}
	p.curScope = p.rootScope
	p.scanner.Init(p.file)
	p.next()

	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
			// panic(e)
		}
	}()
	return p.parseFile()
}
