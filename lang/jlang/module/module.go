package module

import (
	"go/ast"
	"io/ioutil"
	"path"
	"strings"

	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/utils"
)

type Module struct {
	Dir   string
	Name  string
	Files []*file.File

	Symbols map[string]*ast.Node
	Imports []*ast.Node
	Exports map[string]*ast.Node
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
		for _, e := range f.Errors {
			utils.ErrLogger.Println(e)
		}
	}
}

func Load(dir string, name string) (*Module, error) {
	// locate all files in directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	mod := &Module{Name: name, Dir: dir}

	// load all *.j files
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".j") {
			file, err := file.Load(path.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}
			mod.Files = append(mod.Files, file)
		}
	}
	return mod, nil
}
