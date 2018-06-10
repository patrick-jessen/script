package nodes

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/lexer"
	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableDeclNode struct {
	Identifier lexer.Token
	Value      parser.ASTNode
}

func (n VariableDeclNode) String() string {
	return fmt.Sprintf("%v\tidentifier=%v", color.Red("VariableDeclNode"), color.Yellow(n.Identifier.Value))
}
