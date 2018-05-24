package parser

import (
	"fmt"
	"io/ioutil"
)

type Parser struct {
	src []byte
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

// Parse parses the source and returns the AST.
func (p *Parser) Parse() (AstNode, error) {
	_, ast, err := p.parseFunctionCall(0)
	return ast, err
}

func (p *Parser) newError(pos int, format string, args ...interface{}) error {
	return &parserError{
		Message:  fmt.Sprintf(format, args...),
		Parser:   p,
		Position: pos,
	}
}

func (p *Parser) parseFunctionCall(start int) (int, AstNode, error) {
	i := start

	var err error
	var ident AstNode

	i, ident, err = p.parseIdentifier(i)
	if ident == nil {
		return start, nil, p.newError(i, "expected identifier")
	}

	if p.src[i] != '(' {
		return start, nil, p.newError(i, "expected '('")
	}
	i++

	var args AstNode
	if p.src[i] != ')' {
		i, args, err = p.parseArgList(i)
		if err != nil {
			return start, nil, err
		}
	}

	if p.src[i] != ')' {
		return start, nil, p.newError(i, "expected ')'")
	}
	i++

	return i, &FuncCallNode{
		Ident: ident,
		Args:  args,
	}, nil
}

func (p *Parser) parseIdentifier(start int) (int, AstNode, error) {
	i := start

	c := p.src[i]
	if c < 65 || (c > 90 && c < 97) || c > 122 {
		return start, nil, p.newError(i, "identifier must start with a letter")
	}
	i++

	for i < len(p.src) {
		c = p.src[i]
		if c < 48 || (c > 57 && c < 65) || (c > 90 && c < 97) || c > 122 {
			break
		}
		i++
	}

	return i, &IdentifierNode{Name: string(p.src[start:i])}, nil
}

func (p *Parser) parseArg(start int) (int, AstNode, error) {
	i := start
	var arg AstNode
	var err error

	i, arg, err = p.parseString(start)
	if err == nil {
		return i, arg, err
	}

	i, arg, err = p.parseFunctionCall(start)
	if err == nil {
		return i, arg, err
	}

	i, arg, err = p.parseIdentifier(start)
	if err == nil {
		return i, arg, err
	}

	return start, nil, p.newError(i, "expected string funccall or identifier")
}

func (p *Parser) parseArgList(start int) (int, AstNode, error) {
	i := start

	var err error
	var ast ArgListNode
	var lastIsComma bool

	for i < len(p.src) {
		var arg AstNode
		i, arg, err = p.parseArg(i)
		if arg == nil {
			if err != nil {
				return start, nil, err
			}
			break
		}
		lastIsComma = false

		ast.Args = append(ast.Args, arg)

		if p.src[i] != ',' {
			break
		}
		i++
		if p.src[i] != ' ' {
			return start, nil, p.newError(i, "expected space")
		}
		i++
		lastIsComma = true
	}

	if lastIsComma {
		return start, nil, p.newError(i-1, "argument list ends with a comma")
	}

	return i, ast, nil
}

func (p *Parser) parseString(start int) (int, AstNode, error) {
	i := start

	if p.src[i] != '"' {
		return start, nil, p.newError(i, "expected \"")
	}
	i++

	for i < len(p.src) {
		if p.src[i] == '"' {
			return i + 1, &StringNode{
				Value: string(p.src[start+1 : i]),
			}, nil
		}
		i++
	}

	return start, nil, p.newError(i-1, "expected \", got EOF")
}
