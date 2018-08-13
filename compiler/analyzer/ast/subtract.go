package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/color"
)

type Subtract struct {
	LHS   Expression
	RHS   Expression
	OpPos file.Pos
}

func (s *Subtract) Pos() file.Pos {
	return s.LHS.Pos()
}
func (s *Subtract) Type() Type {
	return s.LHS.Type()
}

func (s Subtract) String() string {
	lhs := fmt.Sprintf("  %v", s.LHS)
	rhs := fmt.Sprintf("  %v", s.RHS)

	return fmt.Sprintf(
		"%v\n%v\n%v",
		color.Red("Subtract"),
		strings.Replace(lhs, "\n", "\n  ", -1),
		strings.Replace(rhs, "\n", "\n  ", -1),
	)
}
func (*Subtract) TypeCheck() {
}
