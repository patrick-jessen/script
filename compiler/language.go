package compiler

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/token"
)

type Language interface {
	// Scan scans the next token
	// Should return false if token should be skipped
	Scan(LanguageScanner) token.Token

	// Parse generates an AST
	Parse(LanguageParser) ast.Node
}

type Tokens interface {
	NextToken() token.Token
}

type LanguageScanner interface {
	Consume() bool
	Next() rune

	Token(typ string) token.Token
	TokenVal(typ string, value string) token.Token
	Skip() token.Token
	Error(fmt string, args ...interface{}) token.Token

	StartCapture()
	StopCapture() string
}

type LanguageParser interface {
	Consume() bool
	Next() token.Token
	Is(typ string) bool
	Expect(typ string) // TODO: this should be able to escape parsing somehow
}
