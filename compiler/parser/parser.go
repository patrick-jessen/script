package parser

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/module"
)

type ASTNode interface{}

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

func (p *Parser) error(format string, rest ...interface{}) *module.SourceError {
	it := p.iter
	if p.iter == len(p.tokens) {
		it--
	}
	return p.mod.Error(p.tokens[it].Position, fmt.Sprintf(format, rest...))
}

func (p *Parser) setError(e *module.SourceError, overRule bool) {
	if p.err == nil || (e.Position >= p.err.Position && overRule) || e.Position > p.err.Position {
		p.err = e
	}
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

			p.setError(se, false)
		}
	}()
	return p.grammarMap[gid].fn(p), nil
}

func (p *Parser) Any(gs ...interface{}) []ASTNode {
	var out []ASTNode

	for {
		ast := p.Opt(gs...)
		if ast == nil {
			switch gs[0].(type) {
			case GrammarID:
				p.setError(p.error("expected %v", p.grammarMap[gs[0].(GrammarID)].name), true)
			}
			return out
		}
		out = append(out, ast)
	}
}

func (p *Parser) Opt(alts ...interface{}) (ret ASTNode) {
	defer func() {
		if e := recover(); e != nil {
			ret = nil
		}
	}()
	return p.One(alts...)
}

func (p *Parser) One(alts ...interface{}) (ret ASTNode) {
	for _, a := range alts {
		switch a.(type) {
		case GrammarID:
			ret = p.oneGrammar(a.(GrammarID))
		case lexer.TokenID:
			ret = p.oneToken(a.(lexer.TokenID))
		case func(*Parser) ASTNode:
			ret = a.(func(*Parser) ASTNode)(p)
		default:
			panic("invalid argument")
		}

		if ret != nil {
			return
		}
	}

	switch alts[0].(type) {
	case GrammarID:
		p.setError(p.error("expected %v", p.grammarMap[alts[0].(GrammarID)].name), true)
	}
	panic(p.err)
}

func (p *Parser) oneGrammar(g GrammarID) ASTNode {
	ast, err := p.tryParse(g)
	if err == nil {
		return ast
	}
	return nil
}

func (p *Parser) oneToken(t lexer.TokenID) ASTNode {
	if p.iter == len(p.tokens) {
		panic(p.error("expected %v but reached EOS", p.tokenNames[t]))
	}

	tok := p.tokens[p.iter]
	if tok.TokenID == t {
		p.iter++
		return tok
	}

	p.setError(p.error(
		"expected %v but got %v",
		p.tokenNames[t],
		p.tokenNames[p.tokens[p.iter].TokenID],
	), false)

	return nil
}

func (p *Parser) Run(m module.Module, tokens lexer.TokenStream) ASTNode {
	p.tokens = tokens
	p.mod = m
	ast := p.grammarMap[0].fn(p)

	if p.iter < len(tokens) {
		if p.err != nil {
			panic(p.err.Error())
		}
		panic(p.error("did not parse further"))
	}

	return ast
}
