package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/utils/color"
)

type FunctionDeclNode struct {
	Identifier lexer.Token
}

func (n FunctionDeclNode) String() string {
	return fmt.Sprintf("%v\tidentifier=%v", color.Red("FunctionDeclNode"), color.Yellow(n.Identifier.Value))
}
