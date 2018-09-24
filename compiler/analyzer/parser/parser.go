// Package parser is responsible for performing syntactical analysis on a token stream
package parser

import (
	"github.com/patrick-jessen/script/compiler/analyzer/ast"
	"github.com/patrick-jessen/script/compiler/analyzer/scanner"

	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

// Parser is used for generating AST
type Parser struct {
	scanner scanner.Scanner // scanner which is used for obtaining tokens
	tok     token.Token     // the current token

	rootScope *Scope // top level scope
	curScope  *Scope // current scope

	Unresolved      []*ast.Identifier // unresolved module-local symbols
	Imports         []*ast.Identifier // imported symbols
	ImportedModules []Import          // names of imported modules
}

// New creates a new parser
func New(file *file.File) (p *Parser) {
	rootScope := NewScope(nil)
	p = &Parser{
		rootScope: rootScope,
		curScope:  rootScope,
	}
	p.scanner.Init(file)
	p.next()
	return p
}

// Run runs the parser and returns the AST
func (p *Parser) Run() ast.Node {
	return p.parseFile()
}

func (p *Parser) Symbols() map[string]ast.Declarable {
	return p.rootScope.symbols
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
		d.Pos().MarkError("redeclaration of symbol '%v'. First declared here: (%v)",
			name, sym.Pos().Info().Link(),
		)
		return
	}
	p.curScope.symbols[name] = d
}
func (p *Parser) Resolve(ident *ast.Identifier) {
	mod := ident.Module.Value
	scope := p.curScope

	if len(mod) == 0 {
		// the symbol belongs to this module
		for scope != nil {
			sym, ok := scope.symbols[ident.Symbol.Value]
			if ok {
				ident.Typ = sym.Type()
				ident.Obj = sym.Ident().Obj
				return
			}
			scope = scope.parent
		}
		// cannot be resolved yet
		p.Unresolved = append(p.Unresolved, ident)

	} else {
		// the symbol belongs to another module.
		// assert that the particular module is imported.
		for _, i := range p.ImportedModules {
			if i.Alias.ID != token.Invalid {
				if i.Alias.Value == mod {
					ident.Module.Value = i.Module.Value
					p.Imports = append(p.Imports, ident)
					return
				}
			} else if i.Module.Value == mod {
				p.Imports = append(p.Imports, ident)
				return
			}
		}
		ident.Module.Pos.MarkError("module '%v' not imported", ident.Module.Value)
	}
}

type Import struct {
	Alias  token.Token
	Module token.Token
}

func (p *Parser) importModule(imp Import) {
	idTok := imp.Alias
	if idTok.ID == token.Invalid {
		idTok = imp.Module
	}

	for _, i := range p.ImportedModules {
		t := i.Alias
		if t.ID == token.Invalid {
			t = i.Module
		}

		if idTok.Value == t.Value {
			idTok.Pos.MarkError("duplicate import '%v'", idTok.Value)
		}
	}
	p.ImportedModules = append(p.ImportedModules, imp)
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
		p.tok.Pos.MarkError("expected %v", id.String())
	} else {
		p.next()
	}
}
