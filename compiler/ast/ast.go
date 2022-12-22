package ast

import "github.com/patrick-jessen/script/compiler/file"

type NodeInfo struct {
	Type string
	Pos  file.Pos

	Name    string
	Literal string

	Children []Node
}

type Node interface {
	Info() NodeInfo
}
