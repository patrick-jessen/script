package file

import (
	"fmt"
	"strings"
)

// Pos represents a position in a source file
type Pos struct {
	file  *File // file that the position belongs to
	index int   // index into the source string
}

// Info obtains information regarding the position
func (p Pos) Info() PosInfo {
	var (
		lineNo      = 1    // line of the position (1-indexed)
		colNo       int    // column of the position (1-indexed)
		linePos     int    // start position current line
		nextLinePos int    // start position next line
		line        string // the line containing the position
	)

	// look for the line containing the position
	for _, lp := range p.file.linePos {
		if lp > p.index {
			nextLinePos = lp
			break
		}
		lineNo++
		linePos = lp
	}

	if nextLinePos == 0 {
		// the next line has not been marked yet
		// crop line from source until EOF
		line = p.file.Source[linePos:]

		// crop line at next '\n' (may not be present if EOF)
		if nl := strings.Index(line, "\n"); nl != -1 {
			line = line[:nl]
		}
	} else {
		// the next line has been marked
		// crop from this line to next line
		line = p.file.Source[linePos : nextLinePos-1]
	}

	// calculate column position (1-indexed)
	colNo = p.index - linePos + 1

	return PosInfo{
		File:     p.file.Path,
		LineNo:   lineNo,
		ColumnNo: colNo,
		Line:     line,
	}
}

// MakeError creates a new error at this position
func (p Pos) MakeError(message string) {
	p.file.errors = append(p.file.errors, &fileError{
		Position: p,
		Message:  message,
	})
}

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
