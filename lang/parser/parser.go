package parser

import (
	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/lang/lexer"
	"github.com/patrick-jessen/script/lang/parser/nodes"
)

const (
	Root parser.GrammarID = iota
	DeclStatement
	FunctionDecl
	VariableDecl
)

var Rules = []parser.Rule{
	parser.NewRule(Root, "root",
		func(p *parser.Parser) parser.ASTNode {
			return &nodes.StatementsNode{
				Stmts: p.AnyGrammar(DeclStatement),
			}
		},
	),
	parser.NewRule(DeclStatement, "declaration",
		func(p *parser.Parser) parser.ASTNode {
			stmt := p.OneGrammar(
				FunctionDecl,
				VariableDecl,
			)
			p.OneToken(lexer.NewLine)
			return stmt
		},
	),
	parser.NewRule(VariableDecl, "variable declaration",
		func(p *parser.Parser) parser.ASTNode {
			p.OneToken(lexer.Var)
			ident := p.OneToken(lexer.Identifier)
			p.OneToken(lexer.Equal)
			val := p.OneToken(lexer.String)
			return nodes.VariableDeclNode{
				Identifier: ident,
				Value:      val,
			}
		},
	),
	parser.NewRule(FunctionDecl, "function declaration",
		func(p *parser.Parser) parser.ASTNode {
			p.OneToken(lexer.Func)
			ident := p.OneToken(lexer.Identifier)
			p.OneToken(lexer.ParentStart)
			p.OneToken(lexer.ParentEnd)
			p.OneToken(lexer.CurlStart)
			p.OneToken(lexer.NewLine)
			p.OneToken(lexer.CurlEnd)
			return nodes.FunctionDeclNode{
				Identifier: ident,
			}
		},
	),
}
