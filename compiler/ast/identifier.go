package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type Identifier struct {
	Token token.Token
	Typ   Type
}

func (i *Identifier) Pos() token.Pos {
	return i.Token.Pos
}

func (i Identifier) String() (out string) {
	out = fmt.Sprintf("[%v", color.Yellow(i.Token.Value))
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
