package ir

import (
	"github.com/patrick-jessen/script/utils/color"
)

// LoadI loads an immediate value
type LoadI struct {
	Reg Register
	Val int
}

func (i *LoadI) ColorString() color.String {
	return color.NewString("%v %v  %v",
		color.Yellow("LoadI"), i.Reg.ColorString(), color.Green(i.Val),
	)
}

func (i *LoadI) Execute(vm VM) {
	vm.SetReg(i.Reg, i.Val)
}
