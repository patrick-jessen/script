package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
)

type TokenNode struct {
	Token token.Token
}

func (t TokenNode) String() string {
	return fmt.Sprintf("%v", t.Token)
}