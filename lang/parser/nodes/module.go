package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/utils/color"
)

type ModuleNode struct {
	Statements []parser.ASTNode
}

func (n ModuleNode) String() (out string) {
	for _, s := range n.Statements {
		out += fmt.Sprintf("%v\n", s)
	}
	return
}

////////////////////////////////////////////////////////

func (n *ModuleNode) Analyze(m module.Module) {
	fmt.Println("Analyzing module:", m.Name())

	modScope := newScope(m)
	for _, s := range n.Statements {
		switch s.(type) {
		case *FunctionDeclNode:
			n := s.(*FunctionDeclNode)
			n.ForwardDeclare(modScope)
		case *VariableDeclNode:
			n := s.(*VariableDeclNode)
			n.ForwardDeclare(modScope)
		default:
			panic("unexpected node type")
		}
	}

	/*

		for _, d := range n.Declarations {
			d.ForwardDeclare(modScope)         // also attached astNodes to the symbol
		}

		for _, s := range modScope.Symbols {
			s.Analyze()
		}

	*/
}

type Symbol interface {
	SourceLink() string
}

type Function struct {
	name  lexer.Token
	typ   string
	mod   module.Module
	scope Scope
}

func (f *Function) SourceLink() string {
	return f.mod.PositionInfo(f.name.Position).Link()
}

type Variable struct {
	name lexer.Token
	typ  string
	mod  module.Module
}

func (f *Variable) SourceLink() string {
	return f.mod.PositionInfo(f.name.Position).Link()
}

type Scope struct {
	symbols map[string]Symbol
	mod     module.Module
}

func newScope(mod module.Module) *Scope {
	return &Scope{
		symbols: make(map[string]Symbol),
		mod:     mod,
	}
}

func (s *Scope) DeclareFunction(name lexer.Token, typ string) *Scope {
	if _, ok := s.symbols[name.Value]; ok {
		panic(s.mod.Error(name.Position, fmt.Sprintf(
			"redeclaration of symbol. First defined here: (%v)",
			s.symbols[name.Value].SourceLink())))
	}

	fmt.Printf("function declared: %v - %v\n", color.Green(name.Value), color.Yellow(typ))

	f := &Function{
		name: name,
		typ:  typ,
		mod:  s.mod,
	}

	s.symbols[name.Value] = f
	return &f.scope
}

func (s *Scope) DeclareVariable(name lexer.Token, typ string) {
	if _, ok := s.symbols[name.Value]; ok {
		panic(s.mod.Error(name.Position, fmt.Sprintf(
			"redeclaration of symbol. First defined here: (%v)",
			s.symbols[name.Value].SourceLink())))
	}

	fmt.Printf("variable declared: %v - %v\n", color.Green(name.Value), color.Yellow(typ))

	v := &Variable{
		name: name,
		typ:  typ,
		mod:  s.mod,
	}

	s.symbols[name.Value] = v
}
