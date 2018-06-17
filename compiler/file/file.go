package file

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"text/tabwriter"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type File struct {
	PosMask token.Pos
	Path    string // path to the file
	Source  string // source of the file

	Errors []error
	Tokens []token.Token // tokens of the file (only valid after running lexer)
}

func Load(mask token.Pos, path string) *File {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &File{
		PosMask: mask,
		Path:    path,
		Source:  string(b),
	}
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
	var iter int
	var lines = strings.Split(f.Source, "\n")
	var p = int(pos) & 0x00FFFFFF

	for i, l := range lines {
		if iter+len(l) >= p {
			return PosInfo{
				File:     f.Path,
				LineNo:   i + 1,        // base-1 indexed
				ColumnNo: p - iter + 1, // base-1 indexed
				Line:     l,
			}
		}
		iter += len(l) + 1 // +1 due to \n removed by split
	}
	panic("invalid position")
}

func (f *File) TokensString() string {
	b := bytes.NewBuffer([]byte{})
	tw := tabwriter.NewWriter(b, 0, 0, 2, ' ', 0)

	for i, t := range f.Tokens {
		if i == len(f.Tokens)-1 {
			fmt.Fprint(tw, f.tokenString(t))
		} else {
			fmt.Fprintln(tw, f.tokenString(t))
		}
	}

	tw.Flush()
	return b.String()
}
func (f *File) tokenString(t token.Token) string {
	return fmt.Sprintf(
		"%v\t%v\t%v",
		f.PosInfo(t.Pos).Link(),
		color.Green(t.Name()),
		color.Yellow(t.Value),
	)
}
