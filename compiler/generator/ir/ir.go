package ir

import (
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

func (r Register) ColorString() color.String {
	return color.NewString("reg%v", color.Red(r))
}

type Local int

func (l Local) ColorString() color.String {
	return color.NewString("loc%v", color.Blue(l))
}

type Instruction interface {
	ColorString() color.String
	Execute(VM)
}

type Function struct {
	Name         string
	NumLocals    int
	Instructions []Instruction
}

func (f *Function) ColorString() color.String {
	// out = fmt.Sprintf("%v %v", color.Blue("func"), color.Red(f.Name))

	// if f.NumLocals > 0 {
	// 	out += fmt.Sprintf(" (%v locals)", f.NumLocals)
	// }
	// out += "\n"

	// for i := 0; i < len(f.Instructions); i++ {
	// 	out += "  " + f.Instructions[i].String() + "\n"
	// }
	return color.NewString("")
}
