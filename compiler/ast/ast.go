package ast

import "github.com/patrick-jessen/script/compiler/token"

type Node interface {
	Pos() token.Pos
	String() string
}

type Expression interface {
	Node
	Type() string
}

type Decl interface {
	Node
	Type() string
	Name() string
}

type Ref interface {
	Node
	Type() string
	Name() string
}
