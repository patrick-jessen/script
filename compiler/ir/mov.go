package ir

import (
	"github.com/patrick-jessen/script/utils/color"
)

// Move moves the value of a register to another register
type Move struct {
	Src Register
	Dst Register
}

func (i *Move) ColorString() color.String {
	return color.NewString("%v  %v  %v", color.Yellow("Move"), i.Dst.ColorString(), i.Src.ColorString())
}

func (i *Move) Execute(vm VM) {
	vm.SetReg(i.Dst, vm.Reg(i.Src))
}
