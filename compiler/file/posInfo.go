package file

import (
	"fmt"
	"strings"
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
