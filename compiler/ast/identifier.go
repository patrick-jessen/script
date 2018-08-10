package ast

import (
	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type Object struct {
	Num int
}

type Identifier struct {
	Symbol token.Token
	Module token.Token
	Typ    Type
	Obj    *Object
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
	if i.Typ.IsResolved {
		return color.NewString("%v %v", color.Blue(i.Typ), color.Yellow(i.Name())).String()
	}
	return color.Yellow(i.Name()).String()
}

func (i *Identifier) Type() Type {
	return i.Typ
}
func (*Identifier) TypeCheck(errFn ErrorFunc) {
}

func (i *Identifier) Ident() *Identifier {
	return i
}
