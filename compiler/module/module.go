package module

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/file"
)

type Module struct {
	dir   string
	name  string
	Files []*file.File

	Symbols map[string]ast.Declarable
	Imports []*ast.Identifier
	Exports map[string]ast.Declarable
}

func (m *Module) HasErrors() bool {
	for _, f := range m.Files {
		if f.HasErrors() {
			return true
		}
	}
	return false
}

func (m *Module) PrintErrors() {
	for _, f := range m.Files {
		f.PrintErrors()
	}
}

func Load(dir string, name string) *Module {
	// locate all files in directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	mod := &Module{name: name, dir: dir}

	// load all *.j files
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".j") {
			file := file.Load(path.Join(dir, f.Name()))
			mod.Files = append(mod.Files, file)
		}
	}
	return mod
}

func (m *Module) Name() string {
	return m.name
}

func (m *Module) Dir() string {
	return m.dir
}
