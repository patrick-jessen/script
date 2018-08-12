package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableAssign struct {
	Identifier *Identifier
	Value      Expression
	EqPos      file.Pos
}

func (v *VariableAssign) Pos() file.Pos {
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

func (v *VariableAssign) Name() string {
	return v.Identifier.Name()
}

func (v *VariableAssign) TypeCheck() {
	v.Value.TypeCheck()

	if !v.Identifier.Type().IsCompatible(v.Value.Type()) {
		v.EqPos.MakeError(fmt.Sprintf("cannot assign type %v to %v",
			v.Value.Type(), v.Identifier.Type()))
	}
}
