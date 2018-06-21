package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type Add struct {
	LHS   Expression
	RHS   Expression
	OpPos token.Pos
}

func (a *Add) Pos() token.Pos {
	return a.LHS.Pos()
}

func (a *Add) Type() Type {
	return a.LHS.Type()
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
func (a *Add) TypeCheck(errFn ErrorFunc) {
	lhsTyp := a.LHS.Type()
	rhsTyp := a.RHS.Type()

	if lhsTyp.IsResolved && rhsTyp.IsResolved {
		if lhsTyp.Return != rhsTyp.Return {
			errFn(a.OpPos, fmt.Sprintf("cannot add types %v and %v", lhsTyp, rhsTyp))
		}
	}
}
