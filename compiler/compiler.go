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
			name := u.Token.Value
			sym, ok := symbols[name]
			if !ok {
				mod.Error(u.Pos(), fmt.Sprintf("unresolved symbol '%v'", name))
				continue
			}

			u.Typ = sym.Type()
		}
	}

	for _, v := range symbols {
		fmt.Println(v)
		v.TypeCheck(mod.Error)
	}

	if mod.HasErrors() {
		fmt.Println("-ERRORS--------------------")
		mod.PrintErrors()
	}

	// link module
	// TODO
}
