package ast

import (
	"github.com/patrick-jessen/script/compiler/token"
)

type ErrorFunc func(token.Pos, string)

type Node interface {
	Pos() token.Pos
	String() string
	TypeCheck(ErrorFunc)
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

type Type struct {
	IsResolved bool
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
func (t Type) IsCompatible(other Type) bool {
	if !t.IsResolved || !other.IsResolved {
		return true
	}
	if t.IsFunction != other.IsFunction {
		return false
	}
	if !t.IsFunction {
		return t.Return == other.Return
	}

	return false
}
