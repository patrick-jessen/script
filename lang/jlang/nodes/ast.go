package nodes

import "github.com/patrick-jessen/script/compiler/ast"

type Expression interface {
	ast.Node
}

type Declarable interface {
	ast.Node
}
