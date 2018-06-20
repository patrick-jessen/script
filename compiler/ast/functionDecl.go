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

func (f FunctionDecl) String() string {
	block := fmt.Sprintf("  %v", f.Block)

	return fmt.Sprintf(
		"%v identifier=%v\t%v\n%v",
		color.Red("FunctionDecl"),
		f.Identifier,
		color.Blue(f.Type()),
		strings.Replace(block, "\n", "\n  ", -1),
	)
}

func (f *FunctionDecl) Type() string {
	typ := "("
	if f.Args != nil {
		for i, t := range f.Args.Types {
			typ += t.Token.Value
			if i < len(f.Args.Types)-1 {
				typ += ","
			}
		}
	}
	typ += ") -> void"
	return typ
}

func (f *FunctionDecl) Pos() token.Pos {
	return f.Identifier.Pos()
}
