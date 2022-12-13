package ast

import (
	"github.com/patrick-jessen/script/utils/file"
)

type File struct {
	Declarations []Node
}

func (n *File) Pos() file.Pos {
	return n.Declarations[0].Pos()
}

func (n *File) Children() []Node {
	return n.Declarations
}

func (n *File) TypeCheck() {
}
