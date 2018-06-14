package analyzer

import (
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

type Analyzer interface {
	Analyze(module.Module)
}

func Run(mod module.Module, ast parser.ASTNode) {
	a, ok := ast.(Analyzer)
	if !ok {
		panic("root AST node is not analyzable")
	}

	a.Analyze(mod)
}
