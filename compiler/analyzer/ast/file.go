package ast

import (
	"github.com/patrick-jessen/script/utils/file"
)

type File struct {
	Declarations []Node
}

func (m *File) Pos() file.Pos {
	return m.Declarations[0].Pos()
}

func (m File) String(level int) (out string) {
	for _, s := range m.Declarations {
		out += s.String(level)
	}
	return
}

func (*File) TypeCheck() {
}
