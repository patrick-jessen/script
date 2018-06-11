package main

import (
	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/lang"
	"github.com/patrick-jessen/script/lang/module"
	"github.com/patrick-jessen/script/linker"
)

func main() {
	comp := compiler.New(lang.Rules())
	comp.Compile(module.Load("./src"))

	linker.Run()
}
