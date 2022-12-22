package compiler

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/token"
)

type Compiler interface {
	Compile(path string) []byte
}

type ScannerImpl interface {
	Scan(LanguageScanner) token.Token
}

type ParserImpl interface {
	Parse(LanguageParser) ast.Node
}

type LanguageScanner interface {
	Consume() bool
	ConsumeChar(rune) bool

	Next() rune
	NextIs(rune) bool

	Token(typ string) token.Token
	TokenVal(typ string, value string) token.Token
	Skip() token.Token
	Error(fmt string, args ...interface{}) token.Token

	StartCapture()
	StopCapture() string
}

type LanguageParser interface {
	Consume() bool
	ConsumeType(typ string) bool

	Next() token.Token
	NextIs(typ string) bool
}
