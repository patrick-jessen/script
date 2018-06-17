package lexer

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

// Run runs the lexer on a file
func Run(file *file.File) bool {
	source := file.Source
	var iter int
	var subStr string
	var match []string
	var lastErrPos int

outer:
	for iter < len(source) {
		subStr = source[iter:]

		for _, rule := range rules {
			match = rule.regexp.FindStringSubmatch(subStr)
			if len(match) > 0 {

				if !rule.omit {
					t := token.Token{
						ID:  rule.tokenID,
						Pos: token.Pos(iter) | file.PosMask,
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

		if iter-lastErrPos > 1 {
			file.Error(token.Pos(iter)|file.PosMask, "unknown token")
		}
		lastErrPos = iter
		iter++
	}
	return len(file.Errors) > 0
}
