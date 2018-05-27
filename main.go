package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/patrick-jessen/script/src/parser"
)

var (
	regIdent  = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*`)
	regString = regexp.MustCompile(`^[^(!')]*`)
)

func FuncCall(p *parser.Parser) parser.AstNode {
	ident := p.One(Identifier)
	p.One("(")
	args := p.Opt(ArgList)
	p.One(")")

	return map[string]interface{}{
		"Type":      "FunctionCall",
		"Function":  ident,
		"Arguments": args,
	}
}

func ArgList(p *parser.Parser) parser.AstNode {
	first := p.One(String)
	rest := p.Any(ArgListRest)

	args := []parser.AstNode{first}
	args = append(args, rest.([]parser.AstNode)...)

	return args
}
func ArgListRest(p *parser.Parser) parser.AstNode {
	p.One(",")
	p.Opt(" ")
	return p.One(String)
}

func String(p *parser.Parser) parser.AstNode {
	p.Expect("string")

	p.One("'")
	str := p.One(regString)
	p.One("'")
	return str
}

func Identifier(p *parser.Parser) parser.AstNode {
	p.Expect("identifier")
	return p.One(regIdent)
}

func main() {
	p := parser.New("src.j")
	ast := p.One(FuncCall)

	if p.Error() != nil {
		fmt.Println(p.Error())
	}

	out, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(out))
}
