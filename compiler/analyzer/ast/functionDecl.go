package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/utils/file"
)

type FunctionDecl struct {
	Identifier *Identifier
	Args       []*Identifier
	Block      *Block
}

func (f *FunctionDecl) Ident() *Identifier {
	return f.Identifier
}

func (f *FunctionDecl) Name() string {
	return f.Identifier.Name()
}

func (f FunctionDecl) String(level int) (out string) {
	out = f.Identifier.Pos().Info().Link()
	out += strings.Repeat("  ", level)

	out += fmt.Sprintf("%v %v\n",
		color.Red("FunctionDecl"),
		f.Identifier.String(0),
	)

	if len(f.Block.Statements) > 0 {
		out += f.Block.String(level + 1)
	}
	return
}

func (f *FunctionDecl) Init() {
	var args []string
	for _, a := range f.Args {
		args = append(args, a.Type().Return)
	}

	f.Identifier.Typ = Type{
		IsResolved: true,
		IsFunction: true,
		Return:     "void",
		Args:       args,
	}
}

func (f *FunctionDecl) Type() Type {
	return f.Identifier.Type()
}

func (f *FunctionDecl) Pos() file.Pos {
	return f.Identifier.Pos()
}
func (f *FunctionDecl) TypeCheck() {
	f.Block.TypeCheck()
}
