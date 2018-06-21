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
	if !v.Identifier.Type().IsCompatible(v.Value.Type()) {
		errFn(v.Value.Pos(), fmt.Sprintf("cannot assign type %v to %v",
			v.Value.Type(), v.Identifier.Type()))
	}
}
