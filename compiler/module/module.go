package module

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/interfaces"
	"github.com/patrick-jessen/script/compiler/token"
)

type Module struct {
	name   string
	source string
	Files  []*file.File
	Comp   interfaces.Compiler
}

func New(dir string, name string) *Module {
	// locate all files in directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	mod := &Module{
		name: name,
	}
	offset := 0

	// load any *.j file
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".j") {
			file := file.Load(mod, offset, path.Join(dir, f.Name()))
			mod.Files = append(mod.Files, file)

			// append to modules source code
			mod.source += file.Source + "\n"
			// update offset
			offset = len(mod.source)
		}
	}
	return mod
}

func (m *Module) Error(pos token.Pos, message string) error {
	return &sourceError{
		Module:   m,
		Position: pos,
		Message:  message,
	}
}

func (m *Module) PosInfo(pos token.Pos) interfaces.PosInfo {
	var iter int
	var lines []string

	for _, f := range m.Files {
		if iter+len(f.Source) >= int(pos) {
			lines = strings.Split(f.Source, "\n")

			for i, l := range lines {
				if iter+len(l) >= int(pos) {
					return interfaces.PosInfo{
						File:     f.Path,
						LineNo:   i + 1,               // base-1 indexed
						ColumnNo: int(pos) - iter + 1, // base-1 indexed
						Line:     l,
					}
				}
				iter += len(l) + 1 // +1 due to \n removed by split
			}
		}
		iter += len(f.Source) + 1 // +1 due to \n between files
	}
	panic("invalid position")
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Compiler() interfaces.Compiler {
	return m.Comp
}
