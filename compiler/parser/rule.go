package parser

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/parser/nodes"
	"github.com/patrick-jessen/script/compiler/token"
)

const (
	Root GrammarID = iota
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

var rules map[GrammarID]rule

func init() {
	rules = map[GrammarID]rule{
		Root: newRule("root",
			func(p *Parser) ast.Node {
				return &nodes.ModuleNode{
					Statements: p.Any(DeclStatement),
				}
			},
		),

		// Declarations ///////////////////////////////////////////////////////
		DeclStatement: newRule("declaration",
			func(p *Parser) ast.Node {
				stmt := p.One(
					FunctionDecl,
					VariableDecl,
				)
				p.One(token.NewLine)
				return stmt
			},
		),
		VariableDecl: newRule("variable declaration",
			func(p *Parser) ast.Node {
				p.One(token.Var)
				ident := p.One(token.Identifier)
				p.One(token.Equal)
				val := p.One(Expression)
				return &nodes.VariableDeclNode{
					Identifier: ident.(*nodes.TokenNode),
					Value:      val,
				}
			},
		),
		FunctionDecl: newRule("function declaration",
			func(p *Parser) ast.Node {
				p.One(token.Func)
				ident := p.One(token.Identifier)
				p.One(token.ParentStart)
				p.One(token.ParentEnd)
				block := p.One(Block)
				return &nodes.FunctionDeclNode{
					Identifier: ident.(*nodes.TokenNode),
					Block:      block,
				}
			},
		),
		// Statements /////////////////////////////////////////////////////////
		Statement: newRule("statement",
			func(p *Parser) ast.Node {
				stmt := p.One(
					FunctionDecl,
					VariableDecl,
					VariableAssign,
					FunctionCall,
				)
				p.One(token.NewLine)
				return stmt
			},
		),
		VariableAssign: newRule("variable assignment",
			func(p *Parser) ast.Node {
				ident := p.One(token.Identifier)
				p.One(token.Equal)
				val := p.One(Expression)

				return &nodes.VariableAssignNode{
					Identifier: ident,
					Value:      val,
				}
			},
		),
		FunctionCall: newRule("function call",
			func(p *Parser) ast.Node {
				ident := p.One(token.Identifier)
				p.One(token.ParentStart)
				args := p.Opt(FunctionCallArgs)
				p.One(token.ParentEnd)
				return &nodes.FunctionCallNode{
					Identifier: ident,
					Args:       args,
				}
			},
		),
		FunctionCallArgs: newRule("arguments",
			func(p *Parser) ast.Node {
				first := p.One(Expression)
				rest := p.Any(func(p *Parser) ast.Node {
					p.One(token.Comma)
					return p.One(Expression)
				})

				ret := []ast.Node{first}
				ret = append(ret, rest...)

				return &nodes.FunctionCallArgsNode{
					Args: ret,
				}
			},
		),

		// Expressions ////////////////////////////////////////////////////////
		Expression: newRule("expression",
			func(p *Parser) ast.Node {
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
		BaseExpression: newRule("expression",
			func(p *Parser) ast.Node {
				exp := p.One(
					FunctionCall,
					token.String,
					token.Integer,
					token.Float,
					token.Identifier,
				)
				return &nodes.ExpressionNode{
					Expression: exp,
				}
			},
		),
		Add: newRule("add",
			func(p *Parser) ast.Node {
				lhs := p.One(
					Divide,
					Parenthesis,
					Multiply,
					BaseExpression,
				)
				p.One(token.Plus)
				rhs := p.One(Expression)

				return handleMathOrder(lhs, rhs, token.Plus)
			},
		),
		Subtract: newRule("subtract",
			func(p *Parser) ast.Node {
				lhs := p.One(
					Divide,
					Multiply,
					Parenthesis,
					BaseExpression,
				)
				p.One(token.Minus)
				rhs := p.One(Expression)

				return handleMathOrder(lhs, rhs, token.Minus)
			},
		),
		Multiply: newRule("multiply",
			func(p *Parser) ast.Node {
				lhs := p.One(
					Parenthesis,
					BaseExpression,
				)
				p.One(token.Asterisk)
				rhs := p.One(
					Multiply,
					Divide,
					Parenthesis,
					BaseExpression,
				)

				return handleMathOrder(lhs, rhs, token.Asterisk)
			},
		),
		Divide: newRule("divide",
			func(p *Parser) ast.Node {
				lhs := p.One(
					Parenthesis,
					BaseExpression,
				)
				p.One(token.Slash)
				rhs := p.One(
					Divide,
					Multiply,
					Parenthesis,
					BaseExpression,
				)

				return handleMathOrder(lhs, rhs, token.Slash)
			},
		),
		Parenthesis: newRule("parenthesis",
			func(p *Parser) ast.Node {
				p.One(token.ParentStart)
				exp := p.One(Expression)
				p.One(token.ParentEnd)
				return &nodes.ExpressionNode{
					Expression: exp,
					Immutable:  true,
				}
			},
		),

		Block: newRule("block",
			func(p *Parser) ast.Node {
				p.One(token.CurlStart)
				p.One(token.NewLine)
				stmts := p.Any(Statement)
				p.One(token.CurlEnd)
				return &nodes.StatementsNode{
					Stmts: stmts,
				}
			},
		),
	}
}

func handleMathOrder(lhs ast.Node, rhs ast.Node, op token.ID) ast.Node {
	var pnode ast.Node
	node := rhs

	if op == token.Plus || op == token.Minus {
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

	} else if op == token.Asterisk || op == token.Slash {
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

	var newLHS ast.Node

	switch op {
	case token.Plus:
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

	case token.Minus:
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

	case token.Asterisk:
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

	case token.Slash:
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

type GrammarID int
type ParseFunction func(*Parser) ast.Node

type rule struct {
	grammarID GrammarID
	name      string
	fn        ParseFunction
}

// newRule creates a new grammar rule.
func newRule(name string, fn ParseFunction) rule {
	return rule{
		name: name,
		fn:   fn,
	}
}
