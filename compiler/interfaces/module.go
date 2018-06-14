package interfaces

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
)

type PosInfo struct {
	File     string
	LineNo   int
	ColumnNo int
	Line     string
}

func (p PosInfo) Link() string {
	return fmt.Sprintf("%v:%v:%v", p.File, p.LineNo, p.ColumnNo)
}

func (p PosInfo) String() string {
	arrow := strings.Repeat(" ", p.ColumnNo-1) + "^"
	return fmt.Sprintf("%v\n%v", p.Line, arrow)
}

type Module interface {
	Name() string
	PosInfo(pos token.Pos) PosInfo
	Error(pos token.Pos, message string) error
	Compiler() Compiler
}

type Compiler interface {
	TokenName(id token.ID) string
}
