package lang

import (
	"github.com/patrick-jessen/script/compiler"

	"github.com/patrick-jessen/script/lang/lexer"
)

func Rules() compiler.Language {
	return compiler.Language{
		LexerRules: lexer.Rules,
		// ParserRules: parser.Rules,
	}
}
