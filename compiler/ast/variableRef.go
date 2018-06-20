package ast

import "github.com/patrick-jessen/script/compiler/token"

type VariableRef struct {
	Identifier *Identifier
	typ        Type
}

func (v *VariableRef) Name() string {
	return v.Identifier.Token.Value
}

func (v *VariableRef) Pos() token.Pos {
	return v.Identifier.Pos()
}

func (v *VariableRef) String() string {
	return v.Identifier.String()
}

func (v *VariableRef) Type() Type {
	return v.typ
}
func (v *VariableRef) SetType(t Type) {
	v.typ = t
}
