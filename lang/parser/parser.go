package parser

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/lang/lexer"
	"github.com/patrick-jessen/script/lang/parser/nodes"
)

const (
	Root parser.GrammarID = iota
	DeclStatement
	FunctionDecl
	VariableDecl
	Statement
	VariableAssign
	Expression
	BaseExpression
	Add
	Subtract
	Multiply
	Divide
	Parenthesis
	Block
)

var Rules = []parser.Rule{
	parser.NewRule(Root, "root",
		func(p *parser.Parser) parser.ASTNode {
			return &nodes.StatementsNode{
				Stmts: p.AnyGrammar(DeclStatement),
			}
		},
	),

	// Declarations ///////////////////////////////////////////////////////
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
			val := p.OneGrammar(Expression)
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
			block := p.OneGrammar(Block)
			return nodes.FunctionDeclNode{
				Identifier: ident,
				Block:      block,
			}
		},
	),
	// Statements /////////////////////////////////////////////////////////
	parser.NewRule(Statement, "statement",
		func(p *parser.Parser) parser.ASTNode {
			stmt := p.OneGrammar(
				FunctionDecl,
				VariableDecl,
				VariableAssign,
			)
			p.OneToken(lexer.NewLine)
			return stmt
		},
	),
	parser.NewRule(VariableAssign, "variable assignment",
		func(p *parser.Parser) parser.ASTNode {
			ident := p.OneToken(lexer.Identifier)
			p.OneToken(lexer.Equal)
			val := p.OneGrammar(Expression)

			return &nodes.VariableAssignNode{
				Identifier: ident,
				Value:      val,
			}
		},
	),

	// Expressions ////////////////////////////////////////////////////////
	parser.NewRule(Expression, "expression",
		func(p *parser.Parser) parser.ASTNode {
			exp := p.OneGrammar(
				Subtract,
				Add,
				Multiply,
				Divide,
				Parenthesis,
				BaseExpression,
			)
			return &nodes.ExpressionNode{
				Expression: exp,
			}
		},
	),
	parser.NewRule(BaseExpression, "expression",
		func(p *parser.Parser) parser.ASTNode {
			exp := p.OneToken(
				lexer.String,
				lexer.Integer,
				lexer.Float,
			)
			return &nodes.ExpressionNode{
				Expression: exp,
			}
		},
	),
	parser.NewRule(Add, "add",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.OneGrammar(
				Divide,
				Parenthesis,
				Multiply,
				BaseExpression,
			)
			p.OneToken(lexer.Plus)
			rhs := p.OneGrammar(Expression)
			return &nodes.AddNode{
				LHS: lhs,
				RHS: rhs,
			}
		},
	),
	parser.NewRule(Subtract, "subtract",
		func(p *parser.Parser) parser.ASTNode {
			p.Debug()
			lhs := p.OneGrammar(
				Divide,
				Multiply,
				Parenthesis,
				BaseExpression,
			)
			fmt.Println("-SUB--------------------")
			p.Debug()
			// fmt.Println("LHS", lhs)
			// fmt.Println("RHS", rhs)
			p.OneToken(lexer.Minus)
			rhs := p.OneGrammar(Expression)

			return &nodes.SubtractNode{
				LHS: lhs,
				RHS: rhs,
			}
		},
	),
	parser.NewRule(Multiply, "multiply",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.OneGrammar(
				// Divide,
				Parenthesis,
				BaseExpression,
			)
			p.OneToken(lexer.Asterisk)

			rhs := p.OneGrammar(
				Multiply,
				Parenthesis,
				BaseExpression,
			)

			fmt.Println("-MUL--------------------")
			p.Debug()
			fmt.Println("LHS", lhs)
			fmt.Println("RHS", rhs)

			return &nodes.MultiplyNode{
				LHS: lhs,
				RHS: rhs,
			}
		},
	),
	parser.NewRule(Divide, "divide",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.OneGrammar(
				Parenthesis,
				BaseExpression,
			)
			p.OneToken(lexer.Slash)
			fmt.Println("-DIUV--------------------")
			p.Debug()
			fmt.Println("LHS", lhs)
			// fmt.Println("RHS", rhs)
			p.Debug()
			rhs := p.OneGrammar(
				// Divide,
				Multiply,
				Parenthesis,
				BaseExpression,
			)

			return &nodes.DivideNode{
				LHS: lhs,
				RHS: rhs,
			}
		},
	),
	parser.NewRule(Parenthesis, "parenthesis",
		func(p *parser.Parser) parser.ASTNode {
			p.OneToken(lexer.ParentStart)
			exp := p.OneGrammar(Expression)
			fmt.Println("Parent", exp)
			p.OneToken(lexer.ParentEnd)
			return &nodes.ExpressionNode{
				Expression: exp,
			}
		},
	),

	parser.NewRule(Block, "block",
		func(p *parser.Parser) parser.ASTNode {
			p.OneToken(lexer.CurlStart)
			p.OneToken(lexer.NewLine)
			stmts := p.AnyGrammar(Statement)
			p.OneToken(lexer.CurlEnd)
			return &nodes.StatementsNode{
				Stmts: stmts,
			}
		},
	),
}
