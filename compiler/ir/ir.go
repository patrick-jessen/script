package ir

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

type VM interface {
	Reg(Register) int
	SetReg(Register, int)

	Stack(Local) int
	SetStack(Local, int)

	Call(string)
	Data(int) int
}

type Register int

func (r Register) String() string {
	return "reg" + color.Red(fmt.Sprintf("%v", int(r)))
}

type Local int

func (l Local) String() string {
	return "loc" + color.Blue(fmt.Sprintf("%v", int(l)))
}

type Instruction interface {
	String() string
	Execute(VM)
}

type Function struct {
	Name         string
	NumLocals    int
	Instructions []Instruction
}

func (f *Function) String() (out string) {
	out = fmt.Sprintf("%v %v", color.Blue("func"), color.Red(f.Name))

	if f.NumLocals > 0 {
		out += fmt.Sprintf(" (%v locals)", f.NumLocals)
	}
	out += "\n"

	for i := 0; i < len(f.Instructions); i++ {
		out += "  " + f.Instructions[i].String() + "\n"
	}
	return
}
