package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableDecl struct {
	Identifier *Identifier
	Value      Expression
}

func (v *VariableDecl) Name() string {
	return v.Identifier.Token.Value
}
func (v *VariableDecl) Type() Type {
	return v.Value.Type()
}
func (v *VariableDecl) Pos() token.Pos {
	return v.Identifier.Pos()
}

func (v *VariableDecl) String() string {
	val := fmt.Sprintf("  %v", v.Value)

	return fmt.Sprintf(
		"%v %v\n%v",
		color.Red("VariableDecl"),
		v.Identifier,
		strings.Replace(val, "\n", "\n  ", -1),
	)
}

func (v *VariableDecl) TypeCheck(errFn ErrorFunc) {
	v.Value.TypeCheck(errFn)

	if !v.Identifier.Type().IsCompatible(v.Value.Type()) {
		errFn(v.Value.Pos(), fmt.Sprintf("cannot assign type %v to %v", v.Value.Type(), v.Identifier.Type()))
	}
}
