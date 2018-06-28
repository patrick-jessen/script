package main

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/vm"
)

func main() {
	comp := compiler.New("./src")
	prog := comp.Run()

	fmt.Println(prog)
	vm.Run(prog)
}
