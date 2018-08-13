package file

import (
	"io/ioutil"
)

// File represents a single source file
type File struct {
	Path    string   // path of the file
	Source  string   // source of the file
	Errors  []*Error // list of reported errors
	linePos []int    // positions of lines
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

func (f *File) Contents() string {
	return f.Source
}

// MarkLine marks the beginning of a new line.
// p should be the index of a '\n' character.
func (f *File) MarkLine(p int) {
	f.linePos = append(f.linePos, p+1)
}

// NewPos creates a Pos from an index into the file
func (f *File) NewPos(p int) Pos {
	return Pos{file: f, index: p}
}

// HasErrors returns whether the file contains errors.
func (f *File) HasErrors() bool {
	return len(f.Errors) > 0
}
