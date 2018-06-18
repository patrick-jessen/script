package compiler

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

func Run(mod *module.Module) {
	fmt.Printf("Compiling module '%v'\n", mod.Name())

	for _, f := range mod.Files {
		ast := parser.Run(f)
		for _, e := range f.Errors {
			fmt.Println(e)
		}
		fmt.Println(ast)
	}

}
