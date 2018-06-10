package parser

type GrammarID int
type ParseFunction func(*Parser) ASTNode

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
