package parser

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/lexer"
)

type TokenNode struct {
	Token lexer.Token
}

func (t TokenNode) String() string {
	return fmt.Sprintf("%v", t.Token)
}
