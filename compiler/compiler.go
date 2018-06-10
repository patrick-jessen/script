package compiler

import (
	"fmt"
	"time"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

type Language struct {
	LexerRules  []lexer.Rule
	ParserRules []parser.Rule
}

type Compiler struct {
	lexer  lexer.Lexer
	parser parser.Parser
}

func New(lang Language) Compiler {
	tokenNames := make(map[lexer.TokenID]string)
	for _, r := range lang.LexerRules {
		tokenNames[r.TokenID] = r.Name
	}

	return Compiler{
		lexer:  lexer.New(lang.LexerRules, tokenNames),
		parser: parser.New(lang.ParserRules, tokenNames),
	}
}

func (c Compiler) Compile(mod module.Module) {
	var start time.Time

	fmt.Println("LEXING ===================================")
	start = time.Now()
	tokens := c.lexer.Run(mod)
	fmt.Println("time:", time.Since(start))
	fmt.Println(tokens)

	fmt.Println("PARSING ==================================")
	start = time.Now()
	ast := c.parser.Run(mod, tokens)
	fmt.Println("time:", time.Since(start))
	fmt.Println(ast)

	// fmt.Println("ANALYZING ================================")
	// fmt.Println("GENERATING ===============================")
}
