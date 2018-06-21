package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type Divide struct {
	LHS   Expression
	RHS   Expression
	OpPos token.Pos
}

func (d *Divide) Pos() token.Pos {
	return d.LHS.Pos()
}

func (d *Divide) Type() Type {
	return d.LHS.Type()
}

func (d Divide) String() string {
	lhs := fmt.Sprintf("  %v", d.LHS)
	rhs := fmt.Sprintf("  %v", d.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Divide"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}
func (*Divide) TypeCheck(errFn ErrorFunc) {
}
