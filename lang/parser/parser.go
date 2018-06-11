package parser

import (
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
	FunctionCall
	FunctionCallArgs
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
				Stmts: p.Any(DeclStatement),
			}
		},
	),

	// Declarations ///////////////////////////////////////////////////////
	parser.NewRule(DeclStatement, "declaration",
		func(p *parser.Parser) parser.ASTNode {
			stmt := p.One(
				FunctionDecl,
				VariableDecl,
			)
			p.One(lexer.NewLine)
			return stmt
		},
	),
	parser.NewRule(VariableDecl, "variable declaration",
		func(p *parser.Parser) parser.ASTNode {
			p.One(lexer.Var)
			ident := p.One(lexer.Identifier)
			p.One(lexer.Equal)
			val := p.One(Expression)
			return nodes.VariableDeclNode{
				Identifier: ident,
				Value:      val,
			}
		},
	),
	parser.NewRule(FunctionDecl, "function declaration",
		func(p *parser.Parser) parser.ASTNode {
			p.One(lexer.Func)
			ident := p.One(lexer.Identifier)
			p.One(lexer.ParentStart)
			p.One(lexer.ParentEnd)
			block := p.One(Block)
			return nodes.FunctionDeclNode{
				Identifier: ident,
				Block:      block,
			}
		},
	),
	// Statements /////////////////////////////////////////////////////////
	parser.NewRule(Statement, "statement",
		func(p *parser.Parser) parser.ASTNode {
			stmt := p.One(
				FunctionDecl,
				VariableDecl,
				VariableAssign,
				FunctionCall,
			)
			p.One(lexer.NewLine)
			return stmt
		},
	),
	parser.NewRule(VariableAssign, "variable assignment",
		func(p *parser.Parser) parser.ASTNode {
			ident := p.One(lexer.Identifier)
			p.One(lexer.Equal)
			val := p.One(Expression)

			return &nodes.VariableAssignNode{
				Identifier: ident,
				Value:      val,
			}
		},
	),
	parser.NewRule(FunctionCall, "function call",
		func(p *parser.Parser) parser.ASTNode {
			ident := p.One(lexer.Identifier)
			p.One(lexer.ParentStart)
			args := p.Opt(FunctionCallArgs)
			p.One(lexer.ParentEnd)
			return &nodes.FunctionCallNode{
				Identifier: ident,
				Args:       args,
			}
		},
	),
	parser.NewRule(FunctionCallArgs, "arguments",
		func(p *parser.Parser) parser.ASTNode {
			first := p.One(Expression)
			rest := p.Any(func(p *parser.Parser) parser.ASTNode {
				p.One(lexer.Comma)
				return p.One(Expression)
			})

			ret := []parser.ASTNode{first}
			ret = append(ret, rest...)

			return ret
		},
	),

	// Expressions ////////////////////////////////////////////////////////
	parser.NewRule(Expression, "expression",
		func(p *parser.Parser) parser.ASTNode {
			exp := p.One(
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
			exp := p.One(
				FunctionCall,
				lexer.String,
				lexer.Integer,
				lexer.Float,
				lexer.Identifier,
			)
			return &nodes.ExpressionNode{
				Expression: exp,
			}
		},
	),
	parser.NewRule(Add, "add",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.One(
				Divide,
				Parenthesis,
				Multiply,
				BaseExpression,
			)
			p.One(lexer.Plus)
			rhs := p.One(Expression)

			return handleMathOrder(lhs, rhs, lexer.Plus)
		},
	),
	parser.NewRule(Subtract, "subtract",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.One(
				Divide,
				Multiply,
				Parenthesis,
				BaseExpression,
			)
			p.One(lexer.Minus)
			rhs := p.One(Expression)

			return handleMathOrder(lhs, rhs, lexer.Minus)
		},
	),
	parser.NewRule(Multiply, "multiply",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.One(
				Parenthesis,
				BaseExpression,
			)
			p.One(lexer.Asterisk)
			rhs := p.One(
				Multiply,
				Divide,
				Parenthesis,
				BaseExpression,
			)

			return handleMathOrder(lhs, rhs, lexer.Asterisk)
		},
	),
	parser.NewRule(Divide, "divide",
		func(p *parser.Parser) parser.ASTNode {
			lhs := p.One(
				Parenthesis,
				BaseExpression,
			)
			p.One(lexer.Slash)
			rhs := p.One(
				Divide,
				Multiply,
				Parenthesis,
				BaseExpression,
			)

			return handleMathOrder(lhs, rhs, lexer.Slash)
		},
	),
	parser.NewRule(Parenthesis, "parenthesis",
		func(p *parser.Parser) parser.ASTNode {
			p.One(lexer.ParentStart)
			exp := p.One(Expression)
			p.One(lexer.ParentEnd)
			return &nodes.ExpressionNode{
				Expression: exp,
				Immutable:  true,
			}
		},
	),

	parser.NewRule(Block, "block",
		func(p *parser.Parser) parser.ASTNode {
			p.One(lexer.CurlStart)
			p.One(lexer.NewLine)
			stmts := p.Any(Statement)
			p.One(lexer.CurlEnd)
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
