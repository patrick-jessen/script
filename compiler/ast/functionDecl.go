package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type FunctionDecl struct {
	Identifier *Identifier
	Args       *FunctionDeclArgs
	Block      *Block
}

func (f *FunctionDecl) Name() string {
	return f.Identifier.Token.Value
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

	if f.Args != nil {
		for _, t := range f.Args.Types {
			args = append(args, t.Token.Value)
		}
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

func (f *FunctionDecl) Pos() token.Pos {
	return f.Identifier.Pos()
}
func (f *FunctionDecl) TypeCheck(errFn ErrorFunc) {
	f.Block.TypeCheck(errFn)
}
