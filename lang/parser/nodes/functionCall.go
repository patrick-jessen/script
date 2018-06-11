package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/utils/color"
)

type FunctionCallNode struct {
	Identifier parser.ASTNode
	Args       parser.ASTNode
}

func (n FunctionCallNode) String() string {
	args := ""
	if n.Args != nil {
		argArr := n.Args.([]parser.ASTNode)
		args += "\n"
		for i, a := range argArr {
			args += fmt.Sprintf("%v", a)
			if i != len(argArr)-1 {
				args += "\n"
			}
		}
	}

	return fmt.Sprintf(
		"%v identifier=%v%v",
		color.Red("FunctionCall"),
		n.Identifier,
		strings.Replace(args, "\n", "\n  ", -1),
	)
}
