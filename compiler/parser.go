package compiler

import (
	"github.com/patrick-jessen/script/utils/ast"
	"github.com/patrick-jessen/script/utils/token"
)

// Parser is used for generating AST
type Parser struct {
	lang   Language
	tokens Tokens
	token  token.Token // the current token
}

// NewParser creates a new parser
func NewParser(lang Language, tokens Tokens) (p *Parser) {
	p = &Parser{
		tokens: tokens,
		lang:   lang,
	}
	p.Consume()
	return p
}

// Parse runs the parser and returns the AST
func (p *Parser) Parse() ast.Node {
	return p.lang.Parse(p)
}

func (p *Parser) Consume() bool {
	if p.token.Type == token.EOF {
		return false
	}
	p.token = p.tokens.NextToken()
	return true
}

// Next returns the next token
func (p *Parser) Next() token.Token {
	return p.token
}

// Is checks if the next token is of a certain type
func (p *Parser) Is(typ string) bool {
	return p.token.Type == typ
}

// Expect attempts to consume a token of an expected type
func (p *Parser) Expect(typ string) {
	if p.token.Type != typ {
		p.token.Pos.MarkError("expected %s", typ)
	} else {
		p.Consume()
	}
}
