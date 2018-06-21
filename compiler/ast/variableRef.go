package ast

import "github.com/patrick-jessen/script/compiler/token"

type VariableRef struct {
	Identifier *Identifier
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
	return v.Identifier.Type()
}
func (v *VariableRef) SetType(t Type) {
	v.Identifier.Typ = t
}

func (v *VariableRef) TypeCheck(errFn ErrorFunc) {
}
