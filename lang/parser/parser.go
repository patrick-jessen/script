package parser

import (
	"fmt"

	lex "github.com/patrick-jessen/script/compiler/lexer"
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
				Add,
				Subtract,
				Multiply,
				Divide,
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

			var pnode parser.ASTNode = nil
			node := rhs

		loop:
			for {
				switch node.(type) {
				case *nodes.ExpressionNode:
					if node.(*nodes.ExpressionNode).Immutable {
						break loop
					}
					node = node.(*nodes.ExpressionNode).Expression
				case *nodes.AddNode:
					pnode = node
					node = node.(*nodes.AddNode).LHS
				case *nodes.SubtractNode:
					pnode = node
					node = node.(*nodes.SubtractNode).LHS
				default:
					if pnode == nil {
						break loop
					}
					newNode := &nodes.AddNode{
						LHS: lhs,
						RHS: node,
					}
					switch pnode.(type) {
					case *nodes.AddNode:
						pnode.(*nodes.AddNode).LHS = newNode
					case *nodes.SubtractNode:
						pnode.(*nodes.SubtractNode).LHS = newNode
					}
					return rhs
				}
			}

			return &nodes.AddNode{
				LHS: lhs,
				RHS: rhs,
			}
		},
	),
	parser.NewRule(Subtract, "subtract",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.OneGrammar(
				Divide,
				Multiply,
				Parenthesis,
				BaseExpression,
			)
			p.OneToken(lexer.Minus)
			rhs := p.OneGrammar(Expression)

			var pnode parser.ASTNode = nil
			node := rhs

		loop:
			for {
				switch node.(type) {
				case *nodes.ExpressionNode:
					if node.(*nodes.ExpressionNode).Immutable {
						break loop
					}
					node = node.(*nodes.ExpressionNode).Expression
				case *nodes.AddNode:
					pnode = node
					node = node.(*nodes.AddNode).LHS
				case *nodes.SubtractNode:
					pnode = node
					node = node.(*nodes.SubtractNode).LHS
				default:
					if pnode == nil {
						break loop
					}
					newNode := &nodes.SubtractNode{
						LHS: lhs,
						RHS: node,
					}
					switch pnode.(type) {
					case *nodes.AddNode:
						pnode.(*nodes.AddNode).LHS = newNode
					case *nodes.SubtractNode:
						pnode.(*nodes.SubtractNode).LHS = newNode
					}
					return rhs
				}
			}

			return &nodes.SubtractNode{
				LHS: lhs,
				RHS: rhs,
			}
		},
	),
	parser.NewRule(Multiply, "multiply",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.OneGrammar(
				Parenthesis,
				BaseExpression,
			)
			p.OneToken(lexer.Asterisk)
			rhs := p.OneGrammar(
				Multiply,
				Divide,
				Parenthesis,
				BaseExpression,
			)

			var pnode parser.ASTNode = nil
			node := rhs

		loop:
			for {
				switch node.(type) {
				case *nodes.ExpressionNode:
					if node.(*nodes.ExpressionNode).Immutable {
						break loop
					}
					node = node.(*nodes.ExpressionNode).Expression
				case *nodes.MultiplyNode:
					pnode = node
					node = node.(*nodes.MultiplyNode).LHS
				case *nodes.DivideNode:
					pnode = node
					node = node.(*nodes.DivideNode).LHS
				default:
					if pnode == nil {
						break loop
					}
					newNode := &nodes.MultiplyNode{
						LHS: lhs,
						RHS: node,
					}
					switch pnode.(type) {
					case *nodes.MultiplyNode:
						pnode.(*nodes.MultiplyNode).LHS = newNode
					case *nodes.DivideNode:
						pnode.(*nodes.DivideNode).LHS = newNode
					}
					return rhs
				}
			}

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
			rhs := p.OneGrammar(
				Divide,
				Multiply,
				Parenthesis,
				BaseExpression,
			)

			var pnode parser.ASTNode = nil
			node := rhs

		loop:
			for {
				switch node.(type) {
				case *nodes.ExpressionNode:
					if node.(*nodes.ExpressionNode).Immutable {
						break loop
					}
					node = node.(*nodes.ExpressionNode).Expression
				case *nodes.MultiplyNode:
					pnode = node
					node = node.(*nodes.MultiplyNode).LHS
				case *nodes.DivideNode:
					pnode = node
					node = node.(*nodes.DivideNode).LHS
				default:
					if pnode == nil {
						break loop
					}
					newNode := &nodes.DivideNode{
						LHS: lhs,
						RHS: node,
					}
					switch pnode.(type) {
					case *nodes.MultiplyNode:
						pnode.(*nodes.MultiplyNode).LHS = newNode
					case *nodes.DivideNode:
						pnode.(*nodes.DivideNode).LHS = newNode
					}
					return rhs
				}
			}

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
				Immutable:  true,
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

func handleMathOrder(lhs parser.ASTNode, rhs parser.ASTNode, op lex.TokenID) parser.ASTNode {
	var pnode parser.ASTNode
	node := rhs

	if op == lexer.Plus || op == lexer.Minus {
	loopPlusMinus:
		for {
			switch node.(type) {
			case *nodes.ExpressionNode:
				exp := node.(*nodes.ExpressionNode)
				if exp.Immutable {
					break loopPlusMinus
				}
				node = exp.Expression
			case *nodes.AddNode:
				pnode = node
				node = node.(*nodes.AddNode).LHS
			case *nodes.SubtractNode:
				pnode = node
				node = node.(*nodes.SubtractNode).LHS
			default:
				break loopPlusMinus
			}
		}

	} else if op == lexer.Asterisk || op == lexer.Slash {
	loopMulDiv:
		for {
			switch node.(type) {
			case *nodes.ExpressionNode:
				exp := node.(*nodes.ExpressionNode)
				if exp.Immutable {
					break loopMulDiv
				}
				node = exp.Expression
			case *nodes.MultiplyNode:
				pnode = node
				node = node.(*nodes.MultiplyNode).LHS
			case *nodes.DivideNode:
				pnode = node
				node = node.(*nodes.DivideNode).LHS
			default:
				break loopMulDiv
			}
		}
	} else {
		panic("unexpected operation token")
	}

	var newLHS parser.ASTNode

	switch op {
	case lexer.Plus:
		if pnode == nil {
			return &nodes.AddNode{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &nodes.AddNode{
			LHS: lhs,
			RHS: node,
		}

	case lexer.Minus:
		if pnode == nil {
			return &nodes.SubtractNode{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &nodes.SubtractNode{
			LHS: lhs,
			RHS: node,
		}

	case lexer.Asterisk:
		if pnode == nil {
			return &nodes.MultiplyNode{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &nodes.MultiplyNode{
			LHS: lhs,
			RHS: node,
		}

	case lexer.Slash:
		if pnode == nil {
			return &nodes.DivideNode{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &nodes.DivideNode{
			LHS: lhs,
			RHS: node,
		}
	}

	switch pnode.(type) {
	case *nodes.AddNode:
		pnode.(*nodes.AddNode).LHS = newLHS
	case *nodes.SubtractNode:
		pnode.(*nodes.SubtractNode).LHS = newLHS
	case *nodes.MultiplyNode:
		pnode.(*nodes.MultiplyNode).LHS = newLHS
	case *nodes.DivideNode:
		pnode.(*nodes.DivideNode).LHS = newLHS
	}

	return rhs
}
