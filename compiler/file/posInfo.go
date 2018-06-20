package file

import (
	"fmt"
	"strings"
)

// PosInfo holds information regarding a position
type PosInfo struct {
	File     string // path to the file of the position
	LineNo   int    // line of the position
	ColumnNo int    // column of the position
	Line     string // the line which contains the position
}

// Link returns a link to the position.
// This allows for printing followable links to the console.
// Example output: src/main.j:2:11
func (p PosInfo) Link() string {
	return fmt.Sprintf("%v:%v:%v", p.File, p.LineNo, p.ColumnNo)
}

// String returns the line of the position, with an arrow which points
// to the position.
// Example output:
// func main() {
//      ^
func (p PosInfo) String() string {
	arrow := strings.Repeat(" ", p.ColumnNo-1) + "^"
	return fmt.Sprintf("%v\n%v", p.Line, arrow)
}
