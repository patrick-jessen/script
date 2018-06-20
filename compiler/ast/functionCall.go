package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type FunctionCall struct {
	Identifier *Identifier
	Args       *FunctionCallArgs
	typ        Type
}

func (f *FunctionCall) Name() string {
	return f.Identifier.Token.Value
}

func (f *FunctionCall) Pos() token.Pos {
	return f.Identifier.Pos()
}

func (f FunctionCall) String() (out string) {
	out = fmt.Sprintf(
		"%v identifier=%v\t%v",
		color.Red("FunctionCall"),
		f.Identifier,
		color.Blue(f.Type()),
	)
	if f.Args != nil {
		argArr := f.Args.Args
		args := "\n"
		for i, a := range argArr {
			args += fmt.Sprintf("%v", a)
			if i != len(argArr)-1 {
				args += "\n"
			}
		}
		out += strings.Replace(args, "\n", "\n  ", -1)
	}
	return
}

func (f *FunctionCall) Type() Type {
	return f.typ
}

func (f *FunctionCall) SetType(t Type) {
	f.typ = t
}
