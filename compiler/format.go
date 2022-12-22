package compiler

import (
	"bytes"
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/compiler/utils"
)

func FormatTokens(tokens []token.Token) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 20, 30, 0, '\t', tabwriter.AlignRight)
	for _, t := range tokens {
		w.Write([]byte(fmt.Sprintf("%v\t%v\n", t.Pos.Info().Link(), t)))
	}
	w.Flush()
	return buffer.String()
}

func FormatAST(root ast.Node) string {
	if root == nil {
		return "<none>"
	}

	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 200, 300, 0, '\t', tabwriter.AlignRight)

	type NodeLevel struct {
		Nodes []ast.Node
		Depth int
	}

	stack := []*NodeLevel{{
		Nodes: []ast.Node{root},
		Depth: 0,
	}}

	for {
		if len(stack) == 0 {
			break
		}

		meta := stack[len(stack)-1]
		if len(meta.Nodes) == 0 {
			stack = stack[:len(stack)-1]
			continue
		}

		depth := meta.Depth
		node := meta.Nodes[0]
		meta.Nodes = meta.Nodes[1:]

		nodeInfo := node.Info()

		var info []string

		nodeType := nodeInfo.Type
		info = append(info, utils.Red(nodeType).String())

		// TODO:
		// if n, ok := node.(Typed); ok {
		// 	info = append(info, color.Blue(n.Type().String()).String())
		// }
		if len(nodeInfo.Name) > 0 {
			info = append(info, utils.Yellow(nodeInfo.Name).String())
		}
		if len(nodeInfo.Literal) > 0 {
			info = append(info, utils.Green(nodeInfo.Literal).String())
		}

		w.Write([]byte(
			fmt.Sprintf(
				"%s\t%s%s\n",
				nodeInfo.Pos.Info().Link(),
				strings.Repeat(" ", depth),
				strings.Join(info, " "),
			),
		))

		if len(nodeInfo.Children) > 0 {
			stack = append(stack, &NodeLevel{
				Depth: depth + 1,
				Nodes: nodeInfo.Children,
			})
		}
	}

	w.Flush()
	return buffer.String()
}
