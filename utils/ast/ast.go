package ast

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/utils/file"
)

type Node interface {
	Pos() file.Pos
	Children() []Node

	TypeCheck()
}

type Typed interface {
	Type() Type
}

type Named interface {
	Name() string
}

type Valued interface {
	Value() string
}

type Type struct {
	IsResolved bool
	IsFunction bool
	Return     string
	Args       []string
}

func (t Type) String() (out string) {
	if !t.IsFunction {
		return t.Return
	}
	out = "("
	for i, a := range t.Args {
		out += a
		if i != len(t.Args)-1 {
			out += ","
		}
	}
	out += ") -> " + t.Return
	return
}
func (t Type) IsCompatible(other Type) bool {
	if !t.IsResolved || !other.IsResolved {
		return true
	}
	if t.IsFunction != other.IsFunction {
		return false
	}
	if !t.IsFunction {
		return t.Return == other.Return
	}

	return false
}

func FormatAST(root Node) string {
	var buffer bytes.Buffer
	w := tabwriter.NewWriter(&buffer, 200, 300, 0, '\t', tabwriter.AlignRight)

	type NodeLevel struct {
		Nodes []Node
		Depth int
	}

	stack := []*NodeLevel{{
		Nodes: []Node{root},
		Depth: 0,
	}}

	for {
		// time.Sleep(500 * time.Millisecond)
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

		var info []string

		nodeType := reflect.TypeOf(node).Elem().Name()
		info = append(info, color.Red(nodeType).String())

		if n, ok := node.(Typed); ok {
			info = append(info, color.Blue(n.Type().String()).String())
		}
		if n, ok := node.(Named); ok {
			info = append(info, color.Yellow(n.Name()).String())
		}
		if n, ok := node.(Valued); ok {
			info = append(info, color.Green(n.Value()).String())
		}

		if nodeType != "Block" {
			w.Write([]byte(
				fmt.Sprintf(
					"%s\t%s%s\n",
					node.Pos().Info().Link(),
					strings.Repeat(" ", depth),
					strings.Join(info, " "),
				),
			))
		} else {
			depth--
		}

		if len(node.Children()) > 0 {
			stack = append(stack, &NodeLevel{
				Depth: depth + 1,
				Nodes: node.Children(),
			})
		}

	}

	w.Flush()
	return buffer.String()
}
