package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/file"
)

type File struct {
	Declarations []Node
}

func (m *File) Pos() file.Pos {
	return m.Declarations[0].Pos()
}

func (m File) String() (out string) {
	for _, s := range m.Declarations {
		out += fmt.Sprintf("%v\n", s)
	}
	return
}
func (*File) TypeCheck() {
}
