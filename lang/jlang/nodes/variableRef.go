package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
)

type VariableRef struct {
	Identifier *Identifier
}

func (n *VariableRef) Pos() file.Pos {
	return n.Identifier.Pos()
}

func (n *VariableRef) Children() []ast.Node {
	return nil
}

func (n *VariableRef) Name() string {
	return n.Identifier.Name()
}

func (n *VariableRef) Type() ast.Type {
	return n.Identifier.Type()
}
func (n *VariableRef) SetType(t ast.Type) {
	n.Identifier.Typ = t
}

func (n *VariableRef) TypeCheck() {
}
