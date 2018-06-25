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
	EqPos      token.Pos
}

func (v *VariableAssign) Pos() token.Pos {
	return v.Identifier.Pos()
}

func (v VariableAssign) String() string {
	val := fmt.Sprintf("  %v", v.Value)

	return fmt.Sprintf(
		"%v %v\n%v",
		color.Red("VariableAssign"),
		v.Identifier,
		strings.Replace(val, "\n", "\n  ", -1),
	)
}

func (v *VariableAssign) TypeCheck(errFn ErrorFunc) {
	v.Value.TypeCheck(errFn)

	if !v.Identifier.Type().IsCompatible(v.Value.Type()) {
		errFn(v.EqPos, fmt.Sprintf("cannot assign type %v to %v",
			v.Value.Type(), v.Identifier.Type()))
	}
}
