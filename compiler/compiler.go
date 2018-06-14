package compiler

import (
	"fmt"
	"time"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/compiler/token"
)

type Language struct {
	LexerRules  []lexer.Rule
	ParserRules []parser.Rule
}

type Compiler struct {
	tokenNames map[token.ID]string
	lexer      *lexer.Lexer
	parser     *parser.Parser
}

func New(lang Language) *Compiler {
	tokenNames := make(map[token.ID]string)
	for _, r := range lang.LexerRules {
		if !r.Omit {
			tokenNames[r.TokenID] = r.Name
		}
	}

	return &Compiler{
		lexer:      lexer.New(lang.LexerRules),
		parser:     parser.New(lang.ParserRules, tokenNames),
		tokenNames: tokenNames,
	}
}

func (c *Compiler) Compile(mod *module.Module) {
	fmt.Printf("Compiling module '%v'\n", mod.Name())

	mod.Comp = c

	var start time.Time

	for _, f := range mod.Files {
		fmt.Printf("Compiling file '%v'\n", f.Path)
		fmt.Println("LEXING ===================================")
		start = time.Now()
		err := c.lexer.Run(f)
		if err != nil {
			panic(err)
		}
		fmt.Println("time:", time.Since(start))
		fmt.Println(f.TokensString())

		fmt.Println("PARSING ==================================")
		start = time.Now()
		ast := c.parser.Run(f)
		fmt.Println("time:", time.Since(start))
		fmt.Println(ast)

		// fmt.Println("ANALYZING ================================")
		// start = time.Now()
		// analyzer.Run(mod, ast)
		// fmt.Println("time:", time.Since(start))

		// fmt.Println("GENERATING ===============================")
	}

}

func (c *Compiler) TokenName(id token.ID) string {
	return c.tokenNames[id]
}
