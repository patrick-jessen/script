package nodes

import (
	"github.com/patrick-jessen/script/utils/ast"
)

type Expression interface {
	ast.Node
	ast.Typed
}

type Declarable interface {
	ast.Node
	ast.Typed
	ast.Named
	Ident() *Identifier
}
