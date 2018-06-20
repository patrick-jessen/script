package file

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/patrick-jessen/script/compiler/token"
)

// File represents a single source file
type File struct {
	Path    string    // path to the file
	Source  string    // source of the file
	posMask token.Pos // mask for positions
	errors  []error   // list of reported errors
	linePos []int     // positions of lines
}

// Load loads a source file from disk.
// mask should be unique to the owning module and should be in
// the format: 0xMM000000.
func Load(mask token.Pos, path string) *File {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &File{
		posMask: mask,
		Path:    path,
		Source:  string(b),
	}
}

// MarkLine marks the beginning of a new line.
// p should be the index of a '\n' character.
func (f *File) MarkLine(p int) {
	f.linePos = append(f.linePos, p+1)
}

// Pos creates a token.Pos from an index into the file.
// The posMask is applied to make if possible for modules to trace
// the position back to this file.
func (f *File) Pos(p int) token.Pos {
	return token.Pos(p) | f.posMask
}

// Error registers an error in the file.
// pos should refer to a position in this file.
func (f *File) Error(pos token.Pos, message string) {
	f.errors = append(f.errors, &fileError{
		File:     f,
		Position: pos,
		Message:  message,
	})
}

// HasErrors returns whether the file contains errors.
func (f *File) HasErrors() bool {
	return len(f.errors) > 0
}

// PrintErrors prints file errors to the console.
func (f *File) PrintErrors() {
	for _, e := range f.errors {
		fmt.Println(e)
	}
}

// PosInfo returns the file-related information of a position.
// pos should refer to a position in this file.
func (f *File) PosInfo(pos token.Pos) PosInfo {
	var p = int(pos) & 0x00FFFFFF // ignore mask

	var (
		lineNo      = 1    // line of the position (1-indexed)
		colNo       int    // column of the position (1-indexed)
		linePos     int    // start position current line
		nextLinePos int    // start position next line
		line        string // the line containing the position
	)

	// look for the line containing the position
	for _, lp := range f.linePos {
		if lp > p {
			nextLinePos = lp
			break
		}
		lineNo++
		linePos = lp
	}

	if nextLinePos == 0 {
		// the next line has not been marked yet
		// crop line from source until EOF
		line = f.Source[linePos:]

		// crop line at next '\n' (may not be present if EOF)
		if nl := strings.Index(line, "\n"); nl != -1 {
			line = line[:nl]
		}
	} else {
		// the next line has been marked
		// crop from this line to next line
		line = f.Source[linePos : nextLinePos-1]
	}

	// calculate column position (1-indexed)
	colNo = p - linePos + 1

	return PosInfo{
		File:     f.Path,
		LineNo:   lineNo,
		ColumnNo: colNo,
		Line:     line,
	}
}
