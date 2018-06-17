package nodes

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/utils/color"
)

type VariableDeclNode struct {
	Identifier *TokenNode
	Value      ast.Node
}

func (n VariableDeclNode) String() string {
	val := fmt.Sprintf("  %v", n.Value)

	return fmt.Sprintf(
		"%v identifier=%v\n%v",
		color.Red("VariableDecl"),
		n.Identifier,
		strings.Replace(val, "\n", "\n  ", -1),
	)
}

func (n *VariableDeclNode) Analyze(mod module.Module) {
	// if(n.Type == nil && n.Value == nil) {
	//		panic("cannot determine type")
	// }
	// if(n.Type.??() != n.Value.Type()) {
	//		panic("incompatible types")
	// }
	//
	// type := n.Type.??() || n.Value.Type()
	//
	// // panics if already declared
	// scope.DeclareVariable(
	//		identifier: n.Identifier,
	//		type: type,
	// )
	//
	// if n.Value != nil {
	//		scope.AssignVariable(
	//			identifier: n.Identifier,
	//			value: n.Value,   /// <----------- if e.g. a string, then must register string data
	// 		)
	// }
}

/**
type Reference interface{}
type StackReference struct {}
type ArgReference struct {}
type HeapReference struct {}
type DataReference struct {}
*/

func (n *VariableDeclNode) ForwardDeclare(s *Scope) {
	s.DeclareVariable(n.Identifier.Token, "int")
}
