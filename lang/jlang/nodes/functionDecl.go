package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
)

type FunctionDecl struct {
	Identifier *Identifier
	Args       []*Identifier
	Block      *Block
}

func (n *FunctionDecl) Pos() file.Pos {
	return n.Identifier.Pos()
}

func (n *FunctionDecl) Children() []ast.Node {
	return []ast.Node{n.Block}
}

func (n *FunctionDecl) Ident() *Identifier {
	return n.Identifier
}

func (n *FunctionDecl) Name() string {
	return n.Identifier.Name()
}

func (n *FunctionDecl) Init() {
	var args []string
	for _, a := range n.Args {
		args = append(args, a.Type().Return)
	}

	n.Identifier.Typ = ast.Type{
		IsResolved: true,
		IsFunction: true,
		Return:     "void",
		Args:       args,
	}
}

func (n *FunctionDecl) Type() ast.Type {
	return n.Identifier.Type()
}

func (n *FunctionDecl) TypeCheck() {
	n.Block.TypeCheck()
}
