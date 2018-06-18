package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableDecl struct {
	Identifier *Identifier
	Value      Node
}

func (v *VariableDecl) Pos() token.Pos {
	return v.Identifier.Pos()
}

func (v *VariableDecl) String() string {
	val := fmt.Sprintf("  %v", v.Value)

	return fmt.Sprintf(
		"%v identifier=%v\n%v",
		color.Red("VariableDecl"),
		v.Identifier,
		strings.Replace(val, "\n", "\n  ", -1),
	)
}
