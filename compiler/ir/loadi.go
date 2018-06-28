package ir

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

// LoadI loads an immediate value
type LoadI struct {
	Reg Register
	Val int
}

func (i *LoadI) String() string {
	return fmt.Sprintf("%v\t%v\t%v",
		color.Yellow("LoadI"), i.Reg, color.Green(i.Val),
	)
}

func (i *LoadI) Execute(vm VM) {
	vm.SetReg(i.Reg, i.Val)
}
