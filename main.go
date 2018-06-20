package main

import (
	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/compiler/module"
)

func main() {
	compiler.Run(module.Load("./src", "main"))
	// linker.Run()
}
