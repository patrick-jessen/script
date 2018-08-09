package ir

import (
	"github.com/patrick-jessen/script/utils/color"
)

// Call is a function call
type Call struct {
	Func string
}

func (i *Call) ColorString() color.String {
	return color.NewString("%v  %v", color.Yellow("Call"), color.Red(i.Func))
}

func (i *Call) Execute(vm VM) {
	vm.Call(i.Func)
}
