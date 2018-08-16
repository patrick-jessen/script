package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/utils/file"
)

type VariableDecl struct {
	Identifier *Identifier
	Value      Expression
}

func (v *VariableDecl) Name() string {
	return v.Identifier.Name()
}
func (v *VariableDecl) Type() Type {
	return v.Value.Type()
}
func (v *VariableDecl) Pos() file.Pos {
	return v.Identifier.Pos()
}

func (v *VariableDecl) String(level int) (out string) {
	out = v.Identifier.Pos().Info().Link()
	out += strings.Repeat("  ", level)

	out += fmt.Sprintf(
		"%v %v\n",
		color.Red("VariableDecl"),
		v.Identifier.String(0),
	)
	out += v.Value.String(level + 1)
	return
}

func (v *VariableDecl) TypeCheck() {
	v.Value.TypeCheck()

	if !v.Identifier.Type().IsCompatible(v.Value.Type()) {
		v.Value.Pos().MarkError(fmt.Sprintf("cannot assign type %v to %v", v.Value.Type(), v.Identifier.Type()))
	}
}

func (v *VariableDecl) Ident() *Identifier {
	return v.Identifier
}
