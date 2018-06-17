package parser

import (
	"github.com/patrick-jessen/script/compiler/ast"
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
				return &ast.Module{
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
				return &ast.VariableDecl{
					Identifier: ident.(*ast.Identifier),
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
				return &ast.FunctionDeclNode{
					Identifier: ident.(*ast.Identifier),
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

				return &ast.VariableAssign{
					Identifier: ident.(*ast.Identifier),
					Value:      val.(ast.Expression),
				}
			},
		),
		FunctionCall: newRule("function call",
			func(p *Parser) ast.Node {
				ident := p.One(token.Identifier)
				p.One(token.ParentStart)
				args := p.Opt(FunctionCallArgs)
				p.One(token.ParentEnd)
				return &ast.FunctionCall{
					Identifier: ident.(*ast.Identifier),
					Args:       args.(*ast.FunctionCallArgs),
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

				return &ast.FunctionCallArgs{
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
				return &ast.ExpressionNode{
					Expression: exp.(ast.Expression),
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
				return &ast.ExpressionNode{
					Expression: exp.(ast.Expression),
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

				return handleMathOrder(
					lhs.(ast.Expression),
					rhs.(ast.Expression), token.Plus)
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

				return handleMathOrder(
					lhs.(ast.Expression),
					rhs.(ast.Expression), token.Minus)
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

				return handleMathOrder(
					lhs.(ast.Expression),
					rhs.(ast.Expression), token.Asterisk)
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

				return handleMathOrder(
					lhs.(ast.Expression),
					rhs.(ast.Expression), token.Slash)
			},
		),
		Parenthesis: newRule("parenthesis",
			func(p *Parser) ast.Node {
				p.One(token.ParentStart)
				exp := p.One(Expression)
				p.One(token.ParentEnd)
				return &ast.ExpressionNode{
					Expression: exp.(ast.Expression),
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
				return &ast.Statements{
					Stmts: stmts,
				}
			},
		),
	}
}

func handleMathOrder(lhs ast.Expression, rhs ast.Expression, op token.ID) ast.Node {
	var pnode ast.Node
	node := rhs

	if op == token.Plus || op == token.Minus {
	loopPlusMinus:
		for {
			switch node.(type) {
			case *ast.ExpressionNode:
				exp := node.(*ast.ExpressionNode)
				if exp.Immutable {
					break loopPlusMinus
				}
				node = exp.Expression
			case *ast.Add:
				pnode = node
				node = node.(*ast.Add).LHS
			case *ast.Subtract:
				pnode = node
				node = node.(*ast.Subtract).LHS
			default:
				break loopPlusMinus
			}
		}

	} else if op == token.Asterisk || op == token.Slash {
	loopMulDiv:
		for {
			switch node.(type) {
			case *ast.ExpressionNode:
				exp := node.(*ast.ExpressionNode)
				if exp.Immutable {
					break loopMulDiv
				}
				node = exp.Expression
			case *ast.Multiply:
				pnode = node
				node = node.(*ast.Multiply).LHS
			case *ast.Divide:
				pnode = node
				node = node.(*ast.Divide).LHS
			default:
				break loopMulDiv
			}
		}
	} else {
		panic("unexpected operation token")
	}

	var newLHS ast.Expression

	switch op {
	case token.Plus:
		if pnode == nil {
			return &ast.Add{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &ast.Add{
			LHS: lhs,
			RHS: node,
		}

	case token.Minus:
		if pnode == nil {
			return &ast.Subtract{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &ast.Subtract{
			LHS: lhs,
			RHS: node,
		}

	case token.Asterisk:
		if pnode == nil {
			return &ast.Multiply{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &ast.Multiply{
			LHS: lhs,
			RHS: node,
		}

	case token.Slash:
		if pnode == nil {
			return &ast.Divide{
				LHS: lhs,
				RHS: rhs,
			}
		}
		newLHS = &ast.Divide{
			LHS: lhs,
			RHS: node,
		}
	}

	switch pnode.(type) {
	case *ast.Add:
		pnode.(*ast.Add).LHS = newLHS
	case *ast.Subtract:
		pnode.(*ast.Subtract).LHS = newLHS
	case *ast.Multiply:
		pnode.(*ast.Multiply).LHS = newLHS
	case *ast.Divide:
		pnode.(*ast.Divide).LHS = newLHS
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
