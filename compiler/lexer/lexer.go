package lexer

import (
	"github.com/patrick-jessen/script/compiler/module"
)

type Lexer struct {
	rules      []Rule
	tokenNames map[TokenID]string
}

func New(rules []Rule, tokenNames map[TokenID]string) (l Lexer) {
	l.rules = rules
	l.tokenNames = tokenNames
	return
}

// Run runs the lexer on a module.
// Returns a token stream.
func (l *Lexer) Run(mod module.Module) (tokens TokenStream) {
	source := mod.Source()
	var iter int
	var subStr string
	var match []string

outer:
	for iter < len(source) {
		subStr = source[iter:]

		for _, rule := range l.rules {
			match = rule.regexp.FindStringSubmatch(subStr)
			if len(match) > 0 {

				if !rule.Omit {
					t := Token{
						TokenID:  rule.TokenID,
						Position: iter,
						lexer:    l,
						module:   mod,
					}
					if len(match) > 1 {
						t.Value = match[1]
					}

					tokens = append(tokens, t)
				}

				iter += len(match[0])

				continue outer
			}
		}
		panic(mod.Error(iter, "unknown token"))
	}

	return
}
