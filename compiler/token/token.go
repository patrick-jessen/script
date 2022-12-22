package token

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/utils"
)

const (
	Invalid = "INVALID"
	EOF     = "EOF"
	Skip    = "SKIP"
	Error   = "ERROR"
)

// Token is an actual token which has been read from the input string.
type Token struct {
	Pos   file.Pos // Position of token
	Type  string   // Type of the token
	Value string   // The captured value of the token (not all types of tokens have values)
}

// String returns a pretty formatting of the token
func (t Token) String() string {
	return utils.NewString("%v %v", utils.Green(t.Type), utils.Yellow(t.Value)).String()
}
