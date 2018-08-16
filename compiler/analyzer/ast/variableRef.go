package ast

import "github.com/patrick-jessen/script/utils/file"

type VariableRef struct {
	Identifier *Identifier
}

func (v *VariableRef) Name() string {
	return v.Identifier.Name()
}

func (v *VariableRef) Pos() file.Pos {
	return v.Identifier.Pos()
}
func (v *VariableRef) String(level int) string {
	return v.Identifier.String(level)
}

func (v *VariableRef) Type() Type {
	return v.Identifier.Type()
}
func (v *VariableRef) SetType(t Type) {
	v.Identifier.Typ = t
}

func (v *VariableRef) TypeCheck() {
}
