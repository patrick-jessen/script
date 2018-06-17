package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type ModuleNode struct {
	Statements []ast.Node
}

func (n ModuleNode) String() (out string) {
	for _, s := range n.Statements {
		out += fmt.Sprintf("%v\n", s)
	}
	return
}

type Symbol interface {
	SourceLink() string
}

type Function struct {
	name  token.Token
	typ   string
	file  *file.File
	scope Scope
}

func (f *Function) SourceLink() string {
	return f.file.PosInfo(f.name.Pos).Link()
}

type Variable struct {
	name token.Token
	typ  string
	file *file.File
}

func (f *Variable) SourceLink() string {
	return f.file.PosInfo(f.name.Pos).Link()
}

type Scope struct {
	symbols map[string]Symbol
	file    *file.File
}

func newScope(file *file.File) *Scope {
	return &Scope{
		symbols: make(map[string]Symbol),
		file:    file,
	}
}

func (s *Scope) DeclareFunction(name token.Token, typ string) *Scope {
	if _, ok := s.symbols[name.Value]; ok {
		s.file.Error(name.Pos, fmt.Sprintf(
			"redeclaration of symbol. First defined here: (%v)",
			s.symbols[name.Value].SourceLink()))
		return nil
	}

	fmt.Printf("function declared: %v - %v\n", color.Green(name.Value), color.Yellow(typ))

	f := &Function{
		name: name,
		typ:  typ,
		file: s.file,
	}

	s.symbols[name.Value] = f
	return &f.scope
}

func (s *Scope) DeclareVariable(name token.Token, typ string) {
	if _, ok := s.symbols[name.Value]; ok {
		s.file.Error(name.Pos, fmt.Sprintf(
			"redeclaration of symbol. First defined here: (%v)",
			s.symbols[name.Value].SourceLink()))
		return
	}

	fmt.Printf("variable declared: %v - %v\n", color.Green(name.Value), color.Yellow(typ))

	v := &Variable{
		name: name,
		typ:  typ,
		file: s.file,
	}

	s.symbols[name.Value] = v
}
