package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/color"
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

func (f FunctionDecl) String() (out string) {
	out = fmt.Sprintf(
		"%v %v",
		color.Red("FunctionDecl"),
		f.Identifier,
	)
	if len(f.Block.Statements) > 0 {
		block := fmt.Sprintf("  %v", f.Block)
		out += "\n" + strings.Replace(block, "\n", "\n  ", -1)
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
