package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
)

type File struct {
	Declarations []ast.Node
}

func (n *File) Pos() file.Pos {
	return n.Declarations[0].Pos()
}

func (n *File) Children() []ast.Node {
	return n.Declarations
}

func (n *File) TypeCheck() {
}
