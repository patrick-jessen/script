package module

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

type Module struct {
	name   string
	source string
	Files  []*file.File
}

func Load(dir string, name string) *Module {
	// locate all files in directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	mod := &Module{
		name: name,
	}

	fileIdx := 0

	// load any *.j file
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".j") {
			file := file.Load(token.Pos(fileIdx<<24), path.Join(dir, f.Name()))
			mod.Files = append(mod.Files, file)

			// append to modules source code
			mod.source += file.Source + "\n"
			// update offset
			fileIdx++
		}
	}
	return mod
}

func (m *Module) Name() string {
	return m.name
}
