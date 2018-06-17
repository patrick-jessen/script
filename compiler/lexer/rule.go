package lexer

import (
	"regexp"

	"github.com/patrick-jessen/script/compiler/token"
)

var rules = []rule{
	// Comment
	newOmitRule(`comment`, `\/\/[^\n]*(\n|$)`),

	newRule(token.Var, `var`, `var`),
	newRule(token.Func, `func`, `func`),
	newRule(token.Import, `import`, `import`),
	newRule(token.Return, `return`, `return`),

	// Special
	newOmitRule(`whitespace`, `[ \t]+`),
	newRule(token.NewLine, `new line`, `[\r\n\t ]+`),
	newRule(token.Identifier, `identifier`, `([a-zA-Z_][a-zA-Z_0-9]*)`),
	// Math
	newRule(token.Equal, `=`, `=`),
	newRule(token.Plus, `+`, `\+`),
	newRule(token.Minus, `-`, `-`),
	newRule(token.Asterisk, `*`, `\*`),
	newRule(token.Slash, `/`, `\/`),
	// Symbols
	newRule(token.ParentStart, `(`, `\(`),
	newRule(token.ParentEnd, `)`, `\)`),
	newRule(token.CurlStart, `{`, `{`),
	newRule(token.CurlEnd, `}`, `}`),
	newRule(token.Comma, `,`, `,`),
	// Literals
	newRule(token.Float, `float`, `((?:0|([1-9][0-9]*))\.[0-9]+)`),
	newRule(token.Integer, `integer`, `(0|[1-9][0-9]*)`),
	newRule(token.String, `string`, `"([^"]*)"`),
}

// Rule is a rule for matching tokens.
type rule struct {
	tokenID token.ID
	name    string
	regexp  *regexp.Regexp
	omit    bool
}

// newRule creates a new token rule.
func newRule(tokenID token.ID, name string, regSrc string) rule {
	return rule{
		tokenID: tokenID,
		name:    name,
		regexp:  regexp.MustCompile("^" + regSrc),
	}
}

// newOmitRule creates a new token rule which will be omitted from the
// token stream.
func newOmitRule(name string, regSrc string) rule {
	return rule{
		name:   name,
		regexp: regexp.MustCompile("^" + regSrc),
		omit:   true,
	}
}
