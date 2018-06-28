package ir

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

// Move moves the value of a register to another register
type Move struct {
	Src Register
	Dst Register
}

func (i *Move) String() string {
	return fmt.Sprintf("%v\t%v\t%v", color.Yellow("Move"), i.Dst, i.Src)
}

func (i *Move) Execute(vm VM) {
	vm.SetReg(i.Dst, vm.Reg(i.Src))
}
