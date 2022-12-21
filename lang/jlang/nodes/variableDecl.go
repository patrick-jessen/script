package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
)

type VariableDecl struct {
	Identifier *Identifier
	Value      Expression
}

func (n *VariableDecl) Pos() file.Pos {
	return n.Identifier.Pos()
}

func (n *VariableDecl) Children() []ast.Node {
	return []ast.Node{n.Value}
}

func (n *VariableDecl) Name() string {
	return n.Identifier.Name()
}

func (n *VariableDecl) Type() ast.Type {
	return n.Value.Type()
}

func (n *VariableDecl) TypeCheck() {
	n.Value.TypeCheck()

	if !n.Identifier.Typ.IsResolved {
		n.Identifier.Typ = n.Value.Type()
	}

	if !n.Identifier.Type().IsCompatible(n.Value.Type()) {
		n.Value.Pos().MarkError(fmt.Sprintf("cannot assign type %v to %v", n.Value.Type(), n.Identifier.Type()))
	}
}

func (n *VariableDecl) Ident() *Identifier {
	return n.Identifier
}
