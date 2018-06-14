package lexer

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

// Lexer holds the lexing rules
type Lexer struct {
	rules []Rule
}

// New creates a new lexer
func New(rules []Rule) *Lexer {
	return &Lexer{rules: rules}
}

// Run runs the lexer on a file
func (l *Lexer) Run(file *file.File) error {
	source := file.Source
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
					t := token.Token{
						ID:  rule.TokenID,
						Pos: token.Pos(iter + file.Offset),
					}
					if len(match) > 1 {
						t.Value = match[1]
					}

					file.Tokens = append(file.Tokens, t)
				}

				iter += len(match[0])
				continue outer
			}
		}
		return file.Error(token.Pos(iter+file.Offset), "unknown token")
	}
	return nil
}
