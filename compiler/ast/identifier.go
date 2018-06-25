package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type Identifier struct {
	Symbol token.Token
	Module token.Token
	Typ    Type
}

func (i *Identifier) Name() (out string) {
	if i.Module.ID != token.Invalid {
		out = i.Module.Value + "."
	}
	out += i.Symbol.Value
	return
}

func (i *Identifier) Pos() token.Pos {
	return i.Symbol.Pos
}

func (i Identifier) String() (out string) {
	out = fmt.Sprintf("[%v", color.Yellow(i.Name()))
	if i.Typ.IsResolved {
		out += fmt.Sprintf(" %v", color.Blue(i.Typ))
	}
	out += "]"
	return
}

func (i *Identifier) Type() Type {
	return i.Typ
}
func (*Identifier) TypeCheck(errFn ErrorFunc) {
}
