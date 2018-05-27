package parser

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
)

type ParseFunc func(*Parser) AstNode

type Parser struct {
	src    []byte
	iter   int
	err    error
	potErr error
	expect string
}

// New creates a new parser.
func New(file string) *Parser {
	b, e := ioutil.ReadFile(file)
	if e != nil {
		panic(e)
	}
	p := &Parser{src: b}
	return p
}

func (p *Parser) clone() *Parser {
	return &Parser{
		src:  p.src,
		iter: p.iter,
		err:  nil,
	}
}

func (p *Parser) One(a interface{}) AstNode {
	defer func() { p.expect = "" }()

	switch a.(type) {
	case func(*Parser) AstNode:
		return p.oneFn(a.(func(*Parser) AstNode))
	case string:
		return p.oneStr(a.(string))
	case *regexp.Regexp:
		return p.oneReg(a.(*regexp.Regexp))
	default:
		panic(reflect.TypeOf(a))
	}
}

func (p *Parser) Opt(a interface{}) AstNode {
	defer func() { p.expect = "" }()

	switch a.(type) {
	case func(*Parser) AstNode:
		return p.optFn(a.(func(*Parser) AstNode))
	case string:
		return p.optStr(a.(string))
	case *regexp.Regexp:
		return p.optReg(a.(*regexp.Regexp))
	default:
		panic(reflect.TypeOf(a))
	}
}
func (p *Parser) Any(a interface{}) AstNode {
	defer func() { p.expect = "" }()

	switch a.(type) {
	case func(*Parser) AstNode:
		return p.anyFn(a.(func(*Parser) AstNode))
	case string:
		return p.anyStr(a.(string))
	case *regexp.Regexp:
		return p.anyReg(a.(*regexp.Regexp))
	default:
		panic(reflect.TypeOf(a))
	}
}
func (p *Parser) Expect(str string) {
	p.expect = str
}

func (p *Parser) oneFn(fn ParseFunc) AstNode {
	if p.err != nil {
		return nil
	}
	newp := p.clone()
	res := fn(newp)
	if newp.err != nil {
		p.err = newp.err
		p.potErr = newp.potErr
		return nil
	}
	p.iter = newp.iter
	return res
}
func (p *Parser) optFn(fn ParseFunc) AstNode {
	if p.err != nil {
		return nil
	}
	newp := p.clone()
	res := fn(newp)
	if newp.err != nil {
		p.potErr = newp.err
		return nil
	}
	p.potErr = newp.potErr
	p.iter = newp.iter
	return res
}
func (p *Parser) anyFn(fn ParseFunc) []AstNode {
	if p.err != nil {
		return nil
	}
	newp := p.clone()
	var results []AstNode

	for {
		res := fn(newp)
		if newp.err != nil {
			p.potErr = newp.err
			return results
		}
		results = append(results, res)
		p.iter = newp.iter
	}
}

func (p *Parser) oneStr(str string) string {
	if p.err != nil {
		return ""
	}
	for i := 0; i < len(str); i++ {
		if p.iter == len(p.src) {
			p.error(p.iter, "expected %v", str)
			return ""
		}

		if p.src[p.iter] != str[i] {
			p.error(p.iter, "expected %v", str)
			return ""
		}

		p.iter++
	}
	return str
}
func (p *Parser) optStr(str string) string {
	if p.err != nil {
		return ""
	}
	for i := 0; i < len(str); i++ {
		if p.iter == len(p.src) {
			p.potError(p.iter, "expected %v", str)
			return ""
		}

		if p.src[p.iter] != str[i] {
			p.potError(p.iter, "expected %v", str)
			return ""
		}

		p.iter++
	}
	return str
}
func (p *Parser) anyStr(str string) []string {
	if p.err != nil {
		return nil
	}
	var results []string
	for {
		for i := 0; i < len(str); i++ {
			if p.iter == len(p.src) {
				p.potError(p.iter, "expected %v", str)
				return results
			}

			if p.src[p.iter] != str[i] {
				p.potError(p.iter, "expected %v", str)
				return results
			}

			p.iter++
		}
		results = append(results, str)
	}
}

func (p *Parser) oneReg(r *regexp.Regexp) string {
	if p.err != nil {
		return ""
	}
	match := r.Find(p.src[p.iter:])
	if match == nil {
		p.error(p.iter, "expected %v", r)
		return ""
	}
	p.iter += len(match)
	return string(match)
}
func (p *Parser) optReg(r *regexp.Regexp) string {
	if p.err != nil {
		return ""
	}
	match := r.Find(p.src[p.iter:])
	if match == nil {
		p.potError(p.iter, "expected %v", r)
		return ""
	}
	p.iter += len(match)
	return string(match)
}
func (p *Parser) anyReg(r *regexp.Regexp) []string {
	if p.err != nil {
		return nil
	}
	var results []string
	for {
		match := r.Find(p.src[p.iter:])
		if match == nil {
			p.potError(p.iter, "expected %v", r)
			return results
		}

		results = append(results, string(match))
		p.iter += len(match)
	}
}

func (p *Parser) error(pos int, format string, args ...interface{}) {
	if len(p.expect) > 0 {
		args[0] = p.expect
	}

	p.err = &parserError{
		Message:  fmt.Sprintf(format, args...),
		Parser:   p,
		Position: pos,
	}
}
func (p *Parser) potError(pos int, format string, args ...interface{}) {
	if len(p.expect) > 0 {
		args[0] = p.expect
	}

	if p.potErr != nil {
		if p.potErr.(*parserError).Position > pos {
			return
		}
	}
	p.potErr = &parserError{
		Message:  fmt.Sprintf(format, args...),
		Parser:   p,
		Position: pos,
	}
}

func (p *Parser) Error() error {
	if p.potErr != nil {
		if p.potErr.(*parserError).Position > p.err.(*parserError).Position {
			return p.potErr
		}
	}
	return p.err
}
