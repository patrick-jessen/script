package ast

import (
	"github.com/patrick-jessen/script/compiler/token"
)

type Node interface {
	Pos() token.Pos
	String() string
}

type Expression interface {
	Node
	Type() Type
}

type Declarable interface {
	Node
	Type() Type
	Name() string
}

type Resolvable interface {
	Node
	Type() Type
	Name() string
	SetType(t Type)
}

type Type struct {
	IsFunction bool
	Return     string
	Args       []string
}

func (t Type) String() (out string) {
	if !t.IsFunction {
		return t.Return
	}
	out = "("
	for i, a := range t.Args {
		out += a
		if i != len(t.Args)-1 {
			out += ","
		}
	}
	out += ") -> " + t.Return
	return
}
