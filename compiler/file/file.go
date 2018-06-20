package file

import (
	"io/ioutil"

	"github.com/patrick-jessen/script/compiler/token"
)

type File struct {
	posMask token.Pos
	Path    string // path to the file
	Source  string // source of the file

	Errors  []error
	linePos []int
}

func (f *File) MarkLine(p int) {
	f.linePos = append(f.linePos, p+1)
}

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

func (f *File) Pos(p int) token.Pos {
	return token.Pos(p) | f.posMask
}

func (f *File) Error(pos token.Pos, message string) {
	f.Errors = append(f.Errors, &fileError{
		File:     f,
		Position: pos,
		Message:  message,
	})
}

func (f *File) NewError(pos token.Pos, message string) error {
	return &fileError{
		File:     f,
		Position: pos,
		Message:  message,
	}
}

func (f *File) PosInfo(pos token.Pos) PosInfo {
	var p = int(pos) & 0x00FFFFFF
	var lineNo = 1
	var linePos int
	var nextLinePos int
	var line string

	for _, lp := range f.linePos {
		if lp > p {
			nextLinePos = lp
			break
		}
		lineNo++
		linePos = lp
	}

	if nextLinePos == 0 {
		line = f.Source[linePos:]
	} else {
		line = f.Source[linePos : nextLinePos-1]
	}

	return PosInfo{
		File:     f.Path,
		LineNo:   lineNo,
		ColumnNo: p - linePos + 1,
		Line:     line,
	}
}
