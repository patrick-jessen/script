package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type Add struct {
	LHS Expression
	RHS Expression
	Typ Type
}

func (a *Add) Pos() token.Pos {
	return a.LHS.Pos()
}

func (a *Add) Type() Type {
	return a.Typ
}

func (a Add) String() string {
	lhs := fmt.Sprintf("  %v", a.LHS)
	rhs := fmt.Sprintf("  %v", a.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Add"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}
