package file

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"text/tabwriter"

	"github.com/patrick-jessen/script/compiler/interfaces"
	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type File struct {
	Module interfaces.Module // module that the file belongs to
	Offset int               // offset of positions within module
	Path   string            // path to the file
	Source string            // source of the file

	Tokens []token.Token // tokens of the file (only valid after running lexer)
}

func Load(mod interfaces.Module, offset int, path string) *File {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &File{
		Module: mod,
		Offset: offset,
		Path:   path,
		Source: string(b),
	}
}

func (f *File) Error(pos token.Pos, msg string) error {
	return f.Module.Error(pos, msg)
}

func (f *File) PosInfo(pos token.Pos) interfaces.PosInfo {
	return f.Module.PosInfo(pos)
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
		f.Module.PosInfo(t.Pos).Link(),
		color.Green(f.Module.Compiler().TokenName(t.ID)),
		color.Yellow(t.Value),
	)
}
