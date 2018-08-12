package ast

import (
	"github.com/patrick-jessen/script/compiler/file"
)

type Node interface {
	Pos() file.Pos
	String() string
	TypeCheck()
}

type Expression interface {
	Node
	Type() Type
}

type Declarable interface {
	Node
	Type() Type
	Name() string
	Ident() *Identifier
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
