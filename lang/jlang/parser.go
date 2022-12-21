package jlang

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/lang/jlang/nodes"
	"github.com/patrick-jessen/script/lang/jlang/tokens"
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/token"
)

func (l *JLang) Parse(p compiler.LanguageParser) ast.Node {
	node := &nodes.File{}

	for p.Is(tokens.Import) {
		l.parseImport(p)
	}

	for !p.Is(token.EOF) {
		decl := l.parseDeclaration(p)
		if decl != nil {
			node.Declarations = append(node.Declarations, decl)
		}
	}

	return node
}

func (l *JLang) parseFunctionCallArgs(p compiler.LanguageParser) *nodes.FunctionCallArgs {
	node := &nodes.FunctionCallArgs{}

	for {
		expr := l.parseExpression(p)
		node.Args = append(node.Args, expr)

		if !p.Is(tokens.Comma) {
			break
		}
		p.Consume()
	}

	return node
}

func (l *JLang) parseFunctionCall(p compiler.LanguageParser, ident *nodes.Identifier) *nodes.FunctionCall {
	node := &nodes.FunctionCall{}

	node.Identifier = ident
	p.Expect(tokens.ParentStart)
	if !p.Is(tokens.ParentEnd) {
		node.Args = l.parseFunctionCallArgs(p)
	}
	node.LastParentPos = p.Next().Pos
	p.Expect(tokens.ParentEnd)
	// TODO:
	// p.Resolve(node.Identifier)
	return node
}

func (l *JLang) parseVariableAssign(p compiler.LanguageParser, ident *nodes.Identifier) *nodes.VariableAssign {
	node := &nodes.VariableAssign{}

	node.Identifier = ident
	node.EqPos = p.Next().Pos
	p.Expect(tokens.Equal)
	node.Value = l.parseExpression(p)

	// TODO:
	// p.Resolve(node.Identifier)
	return node
}

func (l *JLang) parseStatement(p compiler.LanguageParser) (n ast.Node) {
	switch p.Next().Type {
	case tokens.Var:
		n = l.parseVariableDecl(p)
	case tokens.Identifier:
		ident := l.parseIdentifier(p)
		if p.Is(tokens.ParentStart) {
			n = l.parseFunctionCall(p, ident)
		} else {
			n = l.parseVariableAssign(p, ident)
		}
	default:
		p.Next().Pos.MarkError("expected statement")
		return nil
	}

	return
}

func (l *JLang) parseBlock(p compiler.LanguageParser) *nodes.Block {
	node := &nodes.Block{}

	p.Expect(tokens.CurlStart)
	for !p.Is(tokens.CurlEnd) {
		stmt := l.parseStatement(p)
		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
		}
	}
	p.Expect(tokens.CurlEnd)
	return node
}

func (l *JLang) parseExpression(p compiler.LanguageParser) nodes.Expression {
	// Shunting-yard algorithm
	var val nodes.Expression
	var out []interface{}
	var ops []token.Token

loop:
	for {
		switch p.Next().Type {
		case tokens.String:
			val = l.parseString(p)
		case tokens.Identifier:
			ident := l.parseIdentifier(p)
			ref := &nodes.VariableRef{Identifier: ident}
			val = ref
			// TODO:
			// p.Resolve(ref.Identifier)
		case tokens.Integer:
			val = l.parseInteger(p)
		case tokens.Float:
			val = l.parseFloat(p)
		case tokens.ParentStart:
			p.Expect(tokens.ParentStart)
			val = l.parseExpression(p)
			p.Expect(tokens.ParentEnd)
		default:
			p.Next().Pos.MarkError("expected expression")
			return nil
		}
		out = append(out, val)

		switch p.Next().Type {
		case tokens.Plus:
			fallthrough
		case tokens.Minus:
			for i := len(ops) - 1; i >= 0; i-- {
				out = append(out, ops[i])
				ops = ops[:i]
			}
			ops = append(ops, p.Next())
			p.Consume()

		case tokens.Asterisk:
			fallthrough
		case tokens.Slash:
			for i := len(ops) - 1; i >= 0; i-- {
				if ops[i].Type == tokens.Asterisk || ops[i].Type == tokens.Slash {
					out = append(out, ops[i])
					ops = ops[:i]
				}
			}
			ops = append(ops, p.Next())
			p.Consume()
		default:
			break loop
		}
	}

	for i := len(ops) - 1; i >= 0; i-- {
		out = append(out, ops[i])
	}

	var valStack []nodes.Expression
	var expr nodes.Expression

	for _, v := range out {
		switch v.(type) {
		case nil:
			return nil
		case token.Token:
			lhs := valStack[len(valStack)-2]
			rhs := valStack[len(valStack)-1]
			valStack = valStack[:len(valStack)-2]

			tok := v.(token.Token)
			switch tok.Type {
			case tokens.Plus:
				expr = &nodes.Add{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			case tokens.Minus:
				expr = &nodes.Subtract{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			case tokens.Asterisk:
				expr = &nodes.Multiply{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			case tokens.Slash:
				expr = &nodes.Divide{
					LHS:   lhs,
					RHS:   rhs,
					OpPos: tok.Pos,
				}
			}

			valStack = append(valStack, expr)
		default:
			valStack = append(valStack, v.(nodes.Expression))
		}
	}

	return valStack[0]
}

func (l *JLang) parseInteger(p compiler.LanguageParser) *nodes.Integer {
	node := &nodes.Integer{Token: p.Next()}
	p.Expect(tokens.Integer)
	return node
}

func (l *JLang) parseFloat(p compiler.LanguageParser) *nodes.Float {
	node := &nodes.Float{Token: p.Next()}
	p.Expect(tokens.Float)
	return node
}

func (l *JLang) parseString(p compiler.LanguageParser) *nodes.String {
	node := &nodes.String{Token: p.Next()}
	p.Expect(tokens.String)
	return node
}

func (l *JLang) parseIdentifier(p compiler.LanguageParser) *nodes.Identifier {
	node := &nodes.Identifier{}

	ident := p.Next()
	p.Expect(tokens.Identifier)

	if p.Is(tokens.Dot) {
		p.Consume()
		node.Module = ident
		node.Symbol = p.Next()
		p.Expect(tokens.Identifier)
	} else {
		node.Symbol = ident
	}

	return node
}

func (l *JLang) parseFunctionDeclArgs(p compiler.LanguageParser) []*nodes.Identifier {
	var ret []*nodes.Identifier

	for {
		ident := l.parseIdentifier(p)
		ident.Obj = &nodes.Object{}
		ident.Typ = ast.Type{
			IsResolved: true,
			Return:     p.Next().Value,
		}
		p.Expect(tokens.Identifier)

		ret = append(ret, ident)

		if !p.Is(tokens.Comma) {
			break
		}
		p.Consume()
	}
	return ret
}

func (l *JLang) parseType(p compiler.LanguageParser) ast.Type {
	// TODO:
	typ := ast.Type{
		IsResolved: true,
	}

	typ.Return = p.Next().Value
	p.Expect(tokens.Identifier)
	return typ
}

func (l *JLang) parseFunctionDecl(p compiler.LanguageParser) *nodes.FunctionDecl {
	obj := &nodes.Object{}
	node := &nodes.FunctionDecl{}

	p.Expect(tokens.Func)
	node.Identifier = l.parseIdentifier(p)
	node.Identifier.Obj = obj
	p.Expect(tokens.ParentStart)
	if !p.Is(tokens.ParentEnd) {
		node.Args = l.parseFunctionDeclArgs(p)
	}
	p.Expect(tokens.ParentEnd)

	// TODO:
	// p.PushScope()
	// for _, a := range node.Args {
	// 	p.Declare(a)
	// }
	node.Block = l.parseBlock(p)
	// p.PopScope()

	// node.Init()

	// p.Declare(node)
	return node
}
func (l *JLang) parseVariableDecl(p compiler.LanguageParser) *nodes.VariableDecl {
	obj := &nodes.Object{}
	node := &nodes.VariableDecl{}

	p.Expect(tokens.Var)
	node.Identifier = l.parseIdentifier(p)
	node.Identifier.Obj = obj
	if !p.Is(tokens.Equal) {
		node.Identifier.Typ = l.parseType(p)
	}
	p.Expect(tokens.Equal)
	node.Value = l.parseExpression(p)

	// TODO:
	// p.Declare(ast)
	return node
}

func (l *JLang) parseDeclaration(p compiler.LanguageParser) (n ast.Node) {
	switch p.Next().Type {
	case tokens.Func:
		return l.parseFunctionDecl(p)
	case tokens.Var:
		return l.parseVariableDecl(p)
	default:
		p.Next().Pos.MarkError("expected declaration")
		return nil
	}
}

func (l *JLang) parseImport(p compiler.LanguageParser) {
	var alias token.Token
	var module token.Token

	p.Expect(tokens.Import)
	if p.Is(tokens.Identifier) {
		alias = p.Next()
		p.Consume()
	}
	module = p.Next()
	p.Expect(tokens.String)
	fmt.Println("Import", alias, module)

	// TODO:
	// p.ImportModule(parser.Import{Alias: alias, Module: module})
}

// rootScope *Scope // top level scope
// curScope  *Scope // current scope

// Unresolved      []ast.Node // unresolved module-local symbols
// Imports         []ast.Node // imported symbols
// ImportedModules []Import   // names of imported modules

// func (p *Parser) Symbols() map[string]ast.Node {
// 	return p.rootScope.symbols
// }

// type Scope struct {
// 	parent  *Scope
// 	symbols map[string]ast.Node
// }

// func NewScope(parent *Scope) *Scope {
// 	return &Scope{
// 		parent:  parent,
// 		symbols: make(map[string]ast.Node),
// 	}
// }

// func (p *Parser) Declare(d ast.Declarable) {
// 	name := d.Name()
// 	sym, ok := p.curScope.symbols[name]
// 	if ok {
// 		d.Pos().MarkError("redeclaration of symbol '%v'. First declared here: (%v)",
// 			name, sym.Pos().Info().Link(),
// 		)
// 		return
// 	}
// 	p.curScope.symbols[name] = d
// }
// func (p *Parser) Resolve(ident *ast.Identifier) {
// 	mod := ident.Module.Value
// 	scope := p.curScope

// 	if len(mod) == 0 {
// 		// the symbol belongs to this module
// 		for scope != nil {
// 			sym, ok := scope.symbols[ident.Symbol.Value]
// 			if ok {
// 				ident.Typ = sym.Type()
// 				ident.Obj = sym.Ident().Obj
// 				return
// 			}
// 			scope = scope.parent
// 		}
// 		// cannot be resolved yet
// 		p.Unresolved = append(p.Unresolved, ident)

// 	} else {
// 		// the symbol belongs to another module.
// 		// assert that the particular module is imported.
// 		for _, i := range p.ImportedModules {
// 			if i.Alias.ID != token.Invalid {
// 				if i.Alias.Value == mod {
// 					ident.Module.Value = i.Module.Value
// 					p.Imports = append(p.Imports, ident)
// 					return
// 				}
// 			} else if i.Module.Value == mod {
// 				p.Imports = append(p.Imports, ident)
// 				return
// 			}
// 		}
// 		ident.Module.Pos.MarkError("module '%v' not imported", ident.Module.Value)
// 	}
// }

// type Import struct {
// 	Alias  token.Token
// 	Module token.Token
// }

// func (p *Parser) ImportModule(imp Import) {
// 	idTok := imp.Alias
// 	if idTok.ID == token.Invalid {
// 		idTok = imp.Module
// 	}

// 	for _, i := range p.ImportedModules {
// 		t := i.Alias
// 		if t.ID == token.Invalid {
// 			t = i.Module
// 		}

// 		if idTok.Value == t.Value {
// 			idTok.Pos.MarkError("duplicate import '%v'", idTok.Value)
// 		}
// 	}
// 	p.ImportedModules = append(p.ImportedModules, imp)
// }

// func (p *Parser) PushScope() {
// 	p.curScope = NewScope(p.curScope)
// }
// func (p *Parser) PopScope() {
// 	p.curScope = p.curScope.parent
// }
