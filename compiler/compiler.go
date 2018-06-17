package compiler

import (
	"fmt"
	"time"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

func Run(mod *module.Module) {
	fmt.Printf("Compiling module '%v'\n", mod.Name())
	var start time.Time

	for _, f := range mod.Files {
		fmt.Printf("Compiling file '%v'\n", f.Path)
		fmt.Println("LEXING ===================================")
		start = time.Now()
		hasErrs := lexer.Run(f)
		fmt.Println("time:", time.Since(start))
		if hasErrs {
			for _, e := range f.Errors {
				fmt.Println(e)
			}
			return
		}
		fmt.Println(f.TokensString())

		fmt.Println("PARSING ==================================")
		start = time.Now()
		ast := parser.Run(f)
		fmt.Println("time:", time.Since(start))
		for _, e := range f.Errors {
			fmt.Println(e)
		}
		fmt.Println(ast)

		// fmt.Println("ANALYZING ================================")
		// start = time.Now()
		// analyzer.Run(mod, ast)
		// fmt.Println("time:", time.Since(start))

		// fmt.Println("GENERATING ===============================")
	}

}
