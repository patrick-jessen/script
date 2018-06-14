package parser

import "github.com/patrick-jessen/script/compiler/ast"

type GrammarID int
type ParseFunction func(*Parser) ast.Node

type Rule struct {
	grammarID GrammarID
	name      string
	fn        ParseFunction
}

// NewRule creates a new grammar rule.
func NewRule(grammarID GrammarID, name string, fn ParseFunction) Rule {
	return Rule{
		grammarID: grammarID,
		name:      name,
		fn:        fn,
	}
}
