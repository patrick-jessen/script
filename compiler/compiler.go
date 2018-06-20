package compiler

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

func Run(mod *module.Module) {
	fmt.Printf("Compiling module '%v'\n", mod.Name())

	var symbols = map[string]ast.Declarable{}
	var parsers []*parser.Parser

	// compile each source file individually
	for _, f := range mod.Files {
		p := parser.New(f)
		parsers = append(parsers, p)
		p.Run()

		for k, v := range p.Symbols() {
			if sym, ok := symbols[k]; ok {
				mod.Error(v.Pos(), fmt.Sprintf(
					"redeclaration of symbol '%v'. First declared here: (%v)",
					k, mod.PosInfo(sym.Pos()).Link(),
				))
			} else {
				symbols[k] = v
			}
		}
	}

	for _, p := range parsers {
		for _, u := range p.Unresolved {
			sym, ok := symbols[u.Ref.Name()]
			if !ok {
				mod.Error(u.Ref.Pos(), fmt.Sprintf("unresolved symbol '%v'", u.Ref.Name()))
				continue
			}

			u.Ref.SetType(sym.Type())
			u.Decl = sym
		}
	}

	if mod.HasErrors() {
		mod.PrintErrors()
	}

	for _, v := range symbols {
		fmt.Println(v)
	}

	// link module
	// TODO
}
