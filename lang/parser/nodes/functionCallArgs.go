package nodes

import (
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

type FunctionCallArgsNode struct {
	Args []parser.ASTNode
}

func (n *FunctionCallArgsNode) Analyze(mod module.Module) {

}
