package ir

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

// LoadD loads the address of data
type LoadD struct {
	Reg  Register
	Data int
}

func (i *LoadD) String() string {
	return fmt.Sprintf("%v\t%v\tdat%v", color.Yellow("LoadD"), i.Reg,
		color.Blue(i.Data))
}

func (i *LoadD) Execute(vm VM) {
	vm.SetReg(i.Reg, vm.Data(i.Data))
}
