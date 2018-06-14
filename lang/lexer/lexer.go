package lexer

import (
	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/token"
)

const (
	Var token.ID = iota
	Func
	Import
	Return
	NewLine
	Identifier
	Equal
	Plus
	Minus
	Asterisk
	Slash
	ParentStart
	ParentEnd
	CurlStart
	CurlEnd
	Comma
	Float
	Integer
	String
)

var Rules = []lexer.Rule{
	// Comment
	lexer.NewOmitRule(`comment`, `\/\/[^\n]*(\n|$)`),

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
	lexer.NewRule(Plus, `+`, `\+`),
	lexer.NewRule(Minus, `-`, `-`),
	lexer.NewRule(Asterisk, `*`, `\*`),
	lexer.NewRule(Slash, `/`, `\/`),
	// Symbols
	lexer.NewRule(ParentStart, `(`, `\(`),
	lexer.NewRule(ParentEnd, `)`, `\)`),
	lexer.NewRule(CurlStart, `{`, `{`),
	lexer.NewRule(CurlEnd, `}`, `}`),
	lexer.NewRule(Comma, `,`, `,`),
	// Literals
	lexer.NewRule(Float, `float`, `((?:0|([1-9][0-9]*))\.[0-9]+)`),
	lexer.NewRule(Integer, `integer`, `(0|[1-9][0-9]*)`),
	lexer.NewRule(String, `string`, `"([^"]*)"`),
}
