package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type Object struct {
	Num int
}

type Identifier struct {
	Symbol token.Token
	Module token.Token
	Typ    ast.Type
	Obj    *Object
}

func (n *Identifier) Pos() file.Pos {
	return n.Symbol.Pos
}

func (n *Identifier) Children() []ast.Node {
	return nil
}

func (n *Identifier) Name() (out string) {
	if n.Module.Type != token.Invalid {
		out = n.Module.Value + "."
	}
	out += n.Symbol.Value
	return
}

func (n *Identifier) Type() ast.Type {
	return n.Typ
}
func (n *Identifier) TypeCheck() {
}

func (n *Identifier) Ident() *Identifier {
	return n
}
