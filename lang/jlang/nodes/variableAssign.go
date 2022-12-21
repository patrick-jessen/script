package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
)

type VariableAssign struct {
	Identifier *Identifier
	Value      Expression
	EqPos      file.Pos
}

func (n *VariableAssign) Pos() file.Pos {
	return n.Identifier.Pos()
}

func (n *VariableAssign) Children() []ast.Node {
	return []ast.Node{n.Value}
}

func (n *VariableAssign) Name() string {
	return n.Identifier.Name()
}

func (n *VariableAssign) TypeCheck() {
	n.Value.TypeCheck()

	if !n.Identifier.Type().IsCompatible(n.Value.Type()) {
		n.EqPos.MarkError(fmt.Sprintf("cannot assign type %v to %v",
			n.Value.Type(), n.Identifier.Type()))
	}
}
