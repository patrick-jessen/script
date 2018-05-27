package parser

type AstNode interface{}

type IdentifierNode struct {
	Name AstNode
}

type FuncCallNode struct {
	Ident AstNode
	Args  AstNode
}

type StringNode struct {
	Value AstNode
}

type ArgListNode struct {
	Args []AstNode
}
