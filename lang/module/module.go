package module

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/patrick-jessen/script/compiler/lexer"
	mod "github.com/patrick-jessen/script/compiler/module"
)

type file struct {
	path   string
	source string
}

func loadFile(path string) *file {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return &file{
		path:   path,
		source: string(b),
	}
}

// Load loads a module given a directory.
func Load(dir string, name string) mod.Module {
	mod := &module{
		name: name,
	}

	// locate all files in directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	// load any *.j file
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".j") {
			file := loadFile(path.Join(dir, f.Name()))
			mod.files = append(mod.files, file)

			// append to modules source code
			mod.source += file.source + "\n"
		}
	}

	return mod
}

type module struct {
	files  []*file
	source string
	tokens lexer.TokenStream
	name   string
}

func (m *module) Name() string {
	return m.name
}

func (m *module) Error(position int, message string) *mod.SourceError {
	return &mod.SourceError{
		Module:   m,
		Position: position,
		Message:  message,
	}
}

func (m *module) Source() string {
	return m.source
}

func (m *module) PositionInfo(pos int) mod.PosInfo {
	var iter int
	var lines []string

	for _, f := range m.files {
		if iter+len(f.source) >= pos {
			lines = strings.Split(f.source, "\n")

			for i, l := range lines {
				if iter+len(l) >= pos {
					return mod.PosInfo{
						File:     f.path,
						LineNo:   i + 1,          // base-1 indexed
						ColumnNo: pos - iter + 1, // base-1 indexed
						Line:     l,
					}
				}
				iter += len(l) + 1 // +1 due to \n removed by split
			}
		}
		iter += len(f.source) + 1 // +1 due to \n between files
	}
	panic("invalid position")
}
