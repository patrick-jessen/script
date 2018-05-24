package parser

type AstNode interface{}

type IdentifierNode struct {
	Name string
}

type FuncCallNode struct {
	Ident AstNode
	Args  AstNode
}

type StringNode struct {
	Value string
}

type ArgListNode struct {
	Args []AstNode
}
