package file

import (
	"fmt"
	"io/ioutil"
)

// File represents a single source file
type File struct {
	Path    string  // path to the file
	Source  string  // source of the file
	errors  []error // list of reported errors
	linePos []int   // positions of lines
}

// Load loads a source file from disk
func Load(path string) *File {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &File{
		Path:   path,
		Source: string(b),
	}
}

// MarkLine marks the beginning of a new line.
// p should be the index of a '\n' character.
func (f *File) MarkLine(p int) {
	f.linePos = append(f.linePos, p+1)
}

// Pos creates a Pos from an index into the file
func (f *File) Pos(p int) Pos {
	return Pos{file: f, index: p}
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
