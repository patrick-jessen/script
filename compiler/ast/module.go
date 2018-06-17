package ast

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
)

type Module struct {
	Statements []Node
}

func (m *Module) Pos() token.Pos {
	return m.Statements[0].Pos()
}

func (m Module) String() (out string) {
	for _, s := range m.Statements {
		out += fmt.Sprintf("%v\n", s)
	}
	return
}
