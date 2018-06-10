package parser

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/module"
)

type Parser struct {
	tokens lexer.TokenStream
	iter   int
	mod    module.Module
	err    *module.SourceError

	grammarMap map[GrammarID]Rule
	tokenNames map[lexer.TokenID]string
}

func New(rules []Rule, tokenNames map[lexer.TokenID]string) (p Parser) {
	p.tokenNames = tokenNames
	p.grammarMap = make(map[GrammarID]Rule)

	for _, r := range rules {
		p.grammarMap[r.grammarID] = r
	}
	return
}

func (p *Parser) error(format string, rest ...interface{}) error {
	it := p.iter
	if p.iter == len(p.tokens) {
		it--
	}
	return p.mod.Error(p.tokens[it].Position, fmt.Sprintf(format, rest...))
}

func (p *Parser) Debug() {
	it := p.iter
	if p.iter == len(p.tokens) {
		it--
	}
	fmt.Println(p.mod.PositionInfo(p.tokens[it].Position).String())
}

func (p *Parser) tryParse(gid GrammarID) (n ASTNode, err error) {
	oldIter := p.iter
	defer func() {
		if e := recover(); e != nil {

			switch e.(type) {
			case *module.SourceError:
			default:
				panic(e)
			}
			se := e.(*module.SourceError)
			err = se
			p.iter = oldIter

			if p.err == nil || se.Position > p.err.Position {
				p.err = se
			}
		}
	}()
	return p.grammarMap[gid].fn(p), nil
}

func (p *Parser) OneGrammar(gs ...GrammarID) ASTNode {
	for _, g := range gs {
		ast, err := p.tryParse(g)
		if err == nil {
			return ast
		}
	}
	return nil
}
func (p *Parser) AnyGrammar(gs ...GrammarID) []ASTNode {
	var out []ASTNode

	for {
		ast := p.OneGrammar(gs...)
		if ast == nil {
			return out
		}
		out = append(out, ast)
	}
}

func (p *Parser) OneToken(ts ...lexer.TokenID) lexer.Token {
	if p.iter == len(p.tokens) {
		panic(p.error("expected %v but reached EOS", p.tokenNames[ts[0]]))
	}

	tok := p.tokens[p.iter]
	for _, t := range ts {
		if tok.TokenID == t {
			p.iter++
			return tok
		}
	}

	panic(p.error(
		"expected %v but got %v",
		p.tokenNames[ts[0]],
		p.tokenNames[p.tokens[p.iter].TokenID],
	))
}

func (p *Parser) Run(m module.Module, tokens lexer.TokenStream) ASTNode {
	p.tokens = tokens
	p.mod = m
	ast := p.grammarMap[0].fn(p)

	if p.iter < len(tokens) {
		if p.err != nil {
			fmt.Println(p.err.Error())
		}
		panic(p.error("did not parse further"))
	}

	return ast
}
