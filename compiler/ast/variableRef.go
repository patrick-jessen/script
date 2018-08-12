package ast

import "github.com/patrick-jessen/script/compiler/file"

type VariableRef struct {
	Identifier *Identifier
}

func (v *VariableRef) Name() string {
	return v.Identifier.Name()
}

func (v *VariableRef) Pos() file.Pos {
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

func (v *VariableRef) TypeCheck() {
}
