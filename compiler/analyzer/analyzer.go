package analyzer

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/patrick-jessen/script/compiler/analyzer/ast"
	"github.com/patrick-jessen/script/compiler/analyzer/parser"

	"github.com/patrick-jessen/script/config"

	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/utils/module"
)

type Analyzer struct {
	workDir    string
	Modules    []*module.Module
	SharedLibs map[string][]string
}

func New(dir string) *Analyzer {
	return &Analyzer{
		workDir: dir,
		Modules: []*module.Module{
			module.Load(dir, "main"),
		},
		SharedLibs: make(map[string][]string),
	}
}

func (a *Analyzer) Run() error {
	var modMap = map[string]*module.Module{}

	// analyze modules
	for i := 0; i < len(a.Modules); i++ {
		mod := a.Modules[i]
		a.analyzeModule(mod)

		modMap[mod.Name()] = mod
	}

	// link modules
	for _, m := range a.Modules {
		for _, i := range m.Imports {
			modName := i.Module.Value
			symName := i.Symbol.Value

			if strings.HasSuffix(modName, ".dll") {
				a.AddExternalSymbol(modName, symName)
				continue
			}

			mod, ok := modMap[modName]
			if !ok {
				continue
			}

			sym, ok := mod.Exports[symName]
			if ok {
				i.Typ = sym.Type()
			} else {
				i.Pos().MarkError(fmt.Sprintf(
					"module '%v' does not export symbol '%v'",
					modName, symName,
				))
			}
		}
	}

	// check types
	a.performTypeCheck(modMap)

	if config.DebugAST {
		l := log.New(os.Stderr, "", 0)
		for k, v := range modMap {
			l.Println(color.NewString("AST for module [%v]:", color.Red(k)))
			for _, sym := range v.Symbols {
				l.Println(ast.FormatAST(sym))
			}
		}
	}

	// print errors (if any)
	if a.hasErrors() {
		a.printErrors()
		return fmt.Errorf("compilation aborted due to errors")
	}
	return nil
}

func (a *Analyzer) AddExternalLib(libName string) {
	if _, ok := a.SharedLibs[libName]; !ok {
		a.SharedLibs[libName] = []string{}
	}
}
func (a *Analyzer) AddExternalSymbol(libName string, symName string) {
	a.SharedLibs[libName] = append(a.SharedLibs[libName], symName)
}

func (a *Analyzer) hasErrors() bool {
	for _, m := range a.Modules {
		if m.HasErrors() {
			return true
		}
	}
	return false
}

func (a *Analyzer) printErrors() {
	for _, m := range a.Modules {
		if m.HasErrors() {
			m.PrintErrors()
		}
	}
}

func (a *Analyzer) performTypeCheck(modMap map[string]*module.Module) {
	defer func() { recover() }()

	for _, mod := range modMap {
		for _, sym := range mod.Symbols {
			sym.TypeCheck()
		}
	}
}

func (a *Analyzer) analyzeModule(mod *module.Module) {
	var symbols = map[string]ast.Declarable{}
	var unresolved []*ast.Identifier
	var imports []*ast.Identifier

	// compile each source file
	for _, f := range mod.Files {
		p := parser.New(f)
		p.Run()

		// keep track of module-scope symbols
		for k, v := range p.Symbols() {
			// detect duplicate declarations
			if sym, ok := symbols[k]; ok {
				v.Pos().MarkError(fmt.Sprintf(
					"redeclaration of symbol '%v'. First declared here: (%v)",
					k, sym.Pos().Info().Link(),
				))
			} else {
				symbols[k] = v
			}
		}

		// keep track of unresolved symbols (module scope)
		unresolved = append(unresolved, p.Unresolved...)

		// keep track of imported symbols
		imports = append(imports, p.Imports...)

		// perform imports
		for _, i := range p.ImportedModules {
			err := a.importModule(i.Module.Value)
			if err != nil {
				i.Module.Pos.MarkError(err.Error())
			}
		}
	}

	// resolve unresolved references within the module
	for _, u := range unresolved {
		name := u.Symbol.Value
		sym, ok := symbols[name]
		if ok {
			u.Typ = sym.Type()
			u.Obj = sym.Ident().Obj
		} else {
			u.Pos().MarkError(fmt.Sprintf("unresolved symbol '%v'", name))
		}
	}

	mod.Symbols = symbols
	mod.Exports = symbols // TODO: not all symbols should be exported
	mod.Imports = imports
}

func (a *Analyzer) importModule(imp string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("module '%v' not found", imp)
		}
	}()

	if strings.HasSuffix(imp, ".dll") {
		a.AddExternalLib(imp)
		return
	}

	path := path.Join(a.workDir, imp)

	// ignore if module is already imported
	for _, m := range a.Modules {
		if m.Dir() == path {
			return
		}
	}
	mod := module.Load(path, imp)
	a.Modules = append(a.Modules, mod)
	return
}
