package compiler

import (
	"fmt"
	"path/filepath"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
)

type Compiler struct {
	workDir string
	modules []*module.Module
}

func New(dir string) *Compiler {
	return &Compiler{
		workDir: dir,
		modules: []*module.Module{
			module.Load(dir, "main"),
		},
	}
}

func (c *Compiler) printErrors() {
	for _, m := range c.modules {
		if m.HasErrors() {
			m.PrintErrors()
		}
	}
}

func (c *Compiler) importModule(imp string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("module '%v' not found", imp)
		}
	}()

	path := filepath.Join(c.workDir, imp)

	// ignore if module is already imported
	for _, m := range c.modules {
		if m.Dir() == path {
			return
		}
	}
	mod := module.Load(path, imp)
	c.modules = append(c.modules, mod)
	return
}

func (c *Compiler) Run() {
	var modMap = map[string]*module.Module{}

	// compile modules
	for i := 0; i < len(c.modules); i++ {
		mod := c.modules[i]
		c.compileModule(mod)

		modMap[mod.Name()] = mod
	}

	// link modules
	for _, m := range c.modules {
		for _, i := range m.Imports {
			modName := i.Module.Value
			symName := i.Symbol.Value

			mod, ok := modMap[modName]
			if !ok {
				continue
			}

			sym, ok := mod.Exports[symName]
			if ok {
				i.Typ = sym.Type()
			} else {
				m.Error(i.Pos(), fmt.Sprintf(
					"module '%v' does not export symbol '%v'",
					modName, symName,
				))
			}
		}
	}

	// perform type check
	for _, mod := range modMap {
		for _, sym := range mod.Symbols {
			sym.TypeCheck(mod.Error)
		}
	}

	// print
	for _, mod := range modMap {
		fmt.Printf("Module '%v' ---------------------\n", mod.Name())
		for _, sym := range mod.Symbols {
			fmt.Println(sym)
		}
	}

	fmt.Println()
	fmt.Println()
	c.printErrors()
}

func (c *Compiler) compileModule(mod *module.Module) {
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
				mod.Error(v.Pos(), fmt.Sprintf(
					"redeclaration of symbol '%v'. First declared here: (%v)",
					k, mod.PosInfo(sym.Pos()).Link(),
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
			err := c.importModule(i.Value)
			if err != nil {
				mod.Error(i.Pos, err.Error())
			}
		}
	}

	// resolve unresolved references within the module
	for _, u := range unresolved {
		name := u.Symbol.Value
		sym, ok := symbols[name]
		if ok {
			u.Typ = sym.Type()
		} else {
			mod.Error(u.Pos(), fmt.Sprintf("unresolved symbol '%v'", name))
		}
	}

	mod.Symbols = symbols
	mod.Exports = symbols // TODO: not all symbols should be exported
	mod.Imports = imports
}
