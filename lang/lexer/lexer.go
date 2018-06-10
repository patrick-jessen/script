package lexer

import "github.com/patrick-jessen/script/compiler/lexer"

const (
	Var lexer.TokenID = iota
	Func
	Import
	Return
	NewLine
	Identifier
	Equal
	ParentStart
	ParentEnd
	CurlStart
	CurlEnd
	Comma
	Integer
	String
)

var Rules = []lexer.Rule{
	lexer.NewRule(Var, `var`, `var`),
	lexer.NewRule(Func, `func`, `func`),
	lexer.NewRule(Import, `import`, `import`),
	lexer.NewRule(Return, `return`, `return`),
	// Special
	lexer.NewOmitRule(`whitespace`, `[ \t]+`),
	lexer.NewRule(NewLine, `new line`, `[\r\n\t ]+`),
	lexer.NewRule(Identifier, `identifier`, `([a-zA-Z_][a-zA-Z_0-9]*)`),
	// Math
	lexer.NewRule(Equal, `=`, `=`),
	// Symbols
	lexer.NewRule(ParentStart, `(`, `\(`),
	lexer.NewRule(ParentEnd, `)`, `\)`),
	lexer.NewRule(CurlStart, `{`, `{`),
	lexer.NewRule(CurlEnd, `}`, `}`),
	lexer.NewRule(Comma, `,`, `,`),
	// Literals
	lexer.NewRule(Integer, `integer`, `(0|[1-9][0-9]*)`),
	lexer.NewRule(String, `string`, `"([^"]*)"`),
}
