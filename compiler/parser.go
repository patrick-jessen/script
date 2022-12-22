package compiler

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/token"
)

// Parser is used for generating AST
type Parser struct {
	lang   ParserImpl
	iter   int
	tokens []token.Token
	token  token.Token // the current token
}

// NewParser creates a new parser
func NewParser(lang ParserImpl, tokens []token.Token) (p *Parser) {
	p = &Parser{
		tokens: tokens,
		lang:   lang,
	}
	p.Reset()
	return p
}

// Sets moves to the beginning of the file
func (p *Parser) Reset() {
	p.iter = -1
	p.Consume()
}

// Parse runs the parser and returns the AST
func (p *Parser) Parse() ast.Node {
	return p.lang.Parse(p)
}

// Consume advances to the next token
func (p *Parser) Consume() bool {
	if p.token.Type == token.EOF {
		return false
	}
	p.iter++
	p.token = p.tokens[p.iter]
	return true
}

// ConsumeType attempts to consume a token of an expected type
func (p *Parser) ConsumeType(typ string) bool {
	if p.token.Type != typ {
		return false
	}
	return p.Consume()
}

// Next returns the next token
func (p *Parser) Next() token.Token {
	return p.token
}

// NextIs checks if the next token is of a certain type
func (p *Parser) NextIs(typ string) bool {
	return p.token.Type == typ
}
