package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/utils/color"
)

type FunctionCall struct {
	Identifier    *Identifier
	Args          *FunctionCallArgs
	LastParentPos file.Pos
}

func (f *FunctionCall) Name() string {
	return f.Identifier.Name()
}

func (f *FunctionCall) Pos() file.Pos {
	return f.Identifier.Pos()
}

func (f FunctionCall) String() (out string) {
	out = fmt.Sprintf(
		"%v %v",
		color.Red("FunctionCall"),
		f.Identifier,
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
	return f.Identifier.Type()
}

func (f *FunctionCall) SetType(t Type) {
	f.Identifier.Typ = t
}
func (f *FunctionCall) TypeCheck() {
	if !f.Type().IsResolved {
		return
	}

	numArgs := 0
	if f.Args != nil {
		f.Args.TypeCheck()
		numArgs = len(f.Args.Args)
	}

	if numArgs != len(f.Type().Args) {
		f.LastParentPos.MakeError(fmt.Sprintf(
			"incorrect number of arguments. Expected %v, got %v",
			len(f.Type().Args), numArgs,
		))
	}
	for i, a := range f.Type().Args {
		if i == numArgs {
			break
		}

		r := f.Args.Args[i].Type().Return
		if r != a {
			f.Args.Args[i].Pos().MakeError(fmt.Sprintf(
				"expected %v, got %v", a, r,
			))
		}
	}
}
