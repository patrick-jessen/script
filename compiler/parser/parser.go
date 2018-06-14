package parser

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

type parseError struct {
	err error
	pos token.Pos
}

type Parser struct {
	iter int
	err  parseError
	file *file.File

	grammarMap map[GrammarID]Rule
	tokenNames map[token.ID]string
}

func New(rules []Rule, tokenNames map[token.ID]string) *Parser {
	p := &Parser{
		grammarMap: make(map[GrammarID]Rule),
		tokenNames: tokenNames,
	}
	for _, r := range rules {
		p.grammarMap[r.grammarID] = r
	}
	return p
}

func (p *Parser) error(format string, rest ...interface{}) parseError {
	it := p.iter
	if p.iter == len(p.file.Tokens) {
		it--
	}
	pos := p.file.Tokens[it].Pos

	return parseError{
		err: p.file.Error(pos, fmt.Sprintf(format, rest...)),
		pos: pos,
	}
}

func (p *Parser) setError(e parseError, overRule bool) {
	if p.err.err == nil || (e.pos >= p.err.pos && overRule) || e.pos > p.err.pos {
		p.err = e
	}
}

func (p *Parser) Debug() {
	it := p.iter
	if p.iter == len(p.file.Tokens) {
		it--
	}
	pos := p.file.Tokens[it].Pos

	fmt.Println(p.file.PosInfo(pos).String())
}

func (p *Parser) tryParse(gid GrammarID) (n ast.Node, err error) {
	oldIter := p.iter
	defer func() {
		if e := recover(); e != nil {

			switch e.(type) {
			case parseError:
			default:
				panic(e)
			}
			se := e.(parseError)
			err = se.err
			p.iter = oldIter

			p.setError(se, false)
		}
	}()
	return p.grammarMap[gid].fn(p), nil
}

func (p *Parser) Any(gs ...interface{}) []ast.Node {
	var out []ast.Node

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

func (p *Parser) Opt(alts ...interface{}) (ret ast.Node) {
	defer func() {
		if e := recover(); e != nil {
			ret = nil
		}
	}()
	return p.One(alts...)
}

func (p *Parser) One(alts ...interface{}) (ret ast.Node) {
	for _, a := range alts {
		switch a.(type) {
		case GrammarID:
			ret = p.oneGrammar(a.(GrammarID))
		case token.ID:
			ret = p.oneToken(a.(token.ID))
		case func(*Parser) ast.Node:
			ret = a.(func(*Parser) ast.Node)(p)
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

func (p *Parser) oneGrammar(g GrammarID) ast.Node {
	ast, err := p.tryParse(g)
	if err == nil {
		return ast
	}
	return nil
}

func (p *Parser) oneToken(t token.ID) ast.Node {
	if p.iter == len(p.file.Tokens) {
		panic(p.error("expected %v but reached EOS", p.tokenNames[t]))
	}

	tok := p.file.Tokens[p.iter]
	if tok.ID == t {
		p.iter++
		return &TokenNode{Token: tok}
	}

	p.setError(p.error(
		"expected %v but got %v",
		p.tokenNames[t],
		p.tokenNames[p.file.Tokens[p.iter].ID],
	), false)

	return nil
}

func (p *Parser) Run(f *file.File) ast.Node {
	p.file = f
	ast := p.grammarMap[0].fn(p)

	if p.iter < len(f.Tokens) {
		if p.err.err != nil {
			panic(p.err.err.Error())
		}
		panic(p.error("did not parse further"))
	}

	return ast
}
