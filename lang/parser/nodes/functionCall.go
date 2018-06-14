package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/module"
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
		argArr := n.Args.(*FunctionCallArgsNode).Args
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

func (n *FunctionCallNode) Analyze(mod module.Module) {
	// will throw if not declared, or if types are not compatible
	//
	// mod.ReferenceFunction(
	//		identifier	string("print"),
	//		type		string("(string,string)->void")
	// )

	// type compatability example:
	//
	// call:
	//		(string, string)->void   // "we dont care about the return value"
	// compatible:
	//		(string, string)->[anything]
	//      (string, string, [opt]...)->[anything]
	//		(string, string...)->[anything]
	//		(string...)->[anything]
	//
	// call:
	//		(int)->int			// "we expect an int return value"
	// compatible:
	//		(int)->int
	//      (int, [opt]...)->int
	//		(int...)->int
	//
	// Note the return value can be anthing when using type inference!
	//
	//
	// A reference to the function must be stored in this struct somehow
}
