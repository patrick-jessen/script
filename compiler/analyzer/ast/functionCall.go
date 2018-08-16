package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/utils/file"
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

func (f FunctionCall) String(level int) (out string) {
	out = f.Identifier.Pos().Info().Link()
	out += strings.Repeat("  ", level)

	out += fmt.Sprintf(
		"%v %v\n",
		color.Red("FunctionCall"),
		f.Identifier.String(0),
	)
	if f.Args != nil {
		argArr := f.Args.Args
		for _, a := range argArr {
			out += a.String(level + 1)
		}
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
		f.LastParentPos.MarkError(fmt.Sprintf(
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
			f.Args.Args[i].Pos().MarkError(fmt.Sprintf(
				"expected %v, got %v", a, r,
			))
		}
	}
}
