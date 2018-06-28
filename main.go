package main

import (
	"github.com/patrick-jessen/script/compiler"
)

func main() {
	comp := compiler.New("./src")
	comp.Run()

	// vm.Run()
}
