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
		"%v identifier=%v\t%v",
		color.Red("FunctionDecl"),
		f.Identifier,
		color.Blue(f.Type()),
	)
	if len(f.Block.Statements) > 0 {
		block := fmt.Sprintf("  %v", f.Block)
		out += "\n" + strings.Replace(block, "\n", "\n  ", -1)
	}
	return
}

func (f *FunctionDecl) Type() Type {
	var args []string

	if f.Args != nil {
		for _, t := range f.Args.Types {
			args = append(args, t.Token.Value)
		}
	}

	return Type{
		IsFunction: true,
		Return:     "void",
		Args:       args,
	}
}

func (f *FunctionDecl) Pos() token.Pos {
	return f.Identifier.Pos()
}
