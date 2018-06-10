package lexer

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/utils/color"
)

// TokenID is a unique identifier for a type of token.
type TokenID int

// Token refers to a token in the source.
type Token struct {
	TokenID  TokenID
	Position int
	Value    string

	lexer  *Lexer
	module module.Module
}

func (t Token) String() string {
	return fmt.Sprintf(
		"%v\t%v %v",
		t.module.PositionInfo(t.Position).Link(),
		color.Green(t.lexer.tokenNames[t.TokenID]),
		color.Yellow(t.Value),
	)
}

// TokenStream is an array of tokens.
type TokenStream []Token

func (ts TokenStream) String() (out string) {
	b := bytes.NewBuffer([]byte{})
	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', 0)

	for i, t := range ts {
		if i == len(ts)-1 {
			fmt.Fprint(tw, t.String())
		} else {
			fmt.Fprintln(tw, t.String())
		}
	}

	tw.Flush()
	return b.String()
}
