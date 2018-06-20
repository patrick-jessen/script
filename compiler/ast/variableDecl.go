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

	typ := v.Value.Type()

	return fmt.Sprintf(
		"%v identifier=%v\t%v\n%v",
		color.Red("VariableDecl"),
		v.Identifier,
		color.Blue(typ),
		strings.Replace(val, "\n", "\n  ", -1),
	)
}
