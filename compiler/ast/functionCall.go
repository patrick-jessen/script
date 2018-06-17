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
}

func (f *FunctionCall) Pos() token.Pos {
	return f.Identifier.Pos()
}

func (f FunctionCall) String() string {
	args := ""
	if f.Args != nil {
		argArr := f.Args.Args
		args += "\n"
		for i, a := range argArr {
			args += fmt.Sprintf("%v", a)
			if i != len(argArr)-1 {
				args += "\n"
			}
		}
	}

	return fmt.Sprintf(
		"%v identifier=%v%v",
		color.Red("FunctionCall"),
		f.Identifier,
		strings.Replace(args, "\n", "\n  ", -1),
	)
}
