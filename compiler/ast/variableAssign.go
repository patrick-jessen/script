package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableAssign struct {
	Identifier *Identifier
	Value      Expression
	typ        Type
}

func (v *VariableAssign) Pos() token.Pos {
	return v.Identifier.Pos()
}

func (v VariableAssign) String() string {
	val := fmt.Sprintf("  %v", v.Value)

	return fmt.Sprintf(
		"%v identifier=%v\t%v\n%v",
		color.Red("VariableAssign"),
		v.Identifier,
		color.Blue(v.Type()),
		strings.Replace(val, "\n", "\n  ", -1),
	)
}

func (v *VariableAssign) Type() Type {
	return v.typ
}
func (v *VariableAssign) SetType(t Type) {
	v.typ = t
}

func (v *VariableAssign) Name() string {
	return v.Identifier.Token.Value
}
