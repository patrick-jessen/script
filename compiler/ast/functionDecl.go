package ast

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type FunctionDecl struct {
	Identifier *Identifier
	Block      *Block
}

func (f FunctionDecl) String() string {
	block := fmt.Sprintf("  %v", f.Block)

	return fmt.Sprintf(
		"%v identifier=%v\n%v",
		color.Red("FunctionDecl"),
		f.Identifier,
		strings.Replace(block, "\n", "\n  ", -1),
	)
}

func (f *FunctionDecl) Pos() token.Pos {
	return f.Identifier.Pos()
}
