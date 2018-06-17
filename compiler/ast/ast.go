package ast

import "github.com/patrick-jessen/script/compiler/token"

type Node interface {
	Pos() token.Pos
}

type Expression interface {
	Node
	Type() string
}
