package ir

import (
	"github.com/patrick-jessen/script/utils/color"
)

// LoadD loads the address of data
type LoadD struct {
	Reg  Register
	Data int
}

func (i *LoadD) ColorString() color.String {
	return color.NewString("%v %v  dat%v", color.Yellow("LoadD"), i.Reg.ColorString(), color.Blue(i.Data))
}

func (i *LoadD) Execute(vm VM) {
	vm.SetReg(i.Reg, vm.Data(i.Data))
}
