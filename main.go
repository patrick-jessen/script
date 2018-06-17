package main

import (
	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/compiler/module"

	"github.com/patrick-jessen/script/linker"
)

func main() {
	compiler.Run(module.Load("./src", "main"))

	linker.Run()
}
