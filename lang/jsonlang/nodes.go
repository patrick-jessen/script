package jsonlang

import (
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

type StringNode struct {
	Value token.Token
}

func (n *StringNode) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:    "string",
		Pos:     n.Value.Pos,
		Literal: n.Value.Value,
	}
}

type NumberNode struct {
	Value token.Token
}

func (n *NumberNode) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:    "number",
		Pos:     n.Value.Pos,
		Literal: n.Value.Value,
	}
}

type BoolNode struct {
	Value token.Token
}

func (n *BoolNode) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:    "bool",
		Pos:     n.Value.Pos,
		Literal: n.Value.Value,
	}
}

type NullNode struct {
	Pos file.Pos
}

func (n *NullNode) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type: "null",
		Pos:  n.Pos,
	}
}

type ArrayNode struct {
	Pos      file.Pos
	Children []ast.Node
}

func (n *ArrayNode) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "array",
		Pos:      n.Pos,
		Children: n.Children,
	}
}

type ObjectNode struct {
	Pos        file.Pos
	Properties []ast.Node
}

func (n *ObjectNode) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "object",
		Pos:      n.Pos,
		Children: n.Properties,
	}
}

type PropertyNode struct {
	Name  token.Token
	Value ast.Node
}

func (n *PropertyNode) Info() ast.NodeInfo {
	return ast.NodeInfo{
		Type:     "property",
		Pos:      n.Name.Pos,
		Name:     n.Name.Value,
		Children: []ast.Node{n.Value},
	}
}
