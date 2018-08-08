package ir

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

// Set sets the value of a stack variable.
type Set struct {
	Var Local
	Reg Register
}

func (i *Set) String() string {
	return fmt.Sprintf("%v   %v  %v",
		color.Yellow("Set"), i.Var, i.Reg)
}

func (i *Set) Execute(vm VM) {
	vm.SetStack(i.Var, vm.Reg(i.Reg))
}
