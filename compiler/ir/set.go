package ir

import (
	"github.com/patrick-jessen/script/utils/color"
)

// Set sets the value of a stack variable.
type Set struct {
	Var Local
	Reg Register
}

func (i *Set) ColorString() color.String {
	return color.NewString("%v   %v  %v",
		color.Yellow("Set"), i.Var.ColorString(), i.Reg.ColorString())
}

func (i *Set) Execute(vm VM) {
	vm.SetStack(i.Var, vm.Reg(i.Reg))
}
