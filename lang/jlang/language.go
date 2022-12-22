package jlang

import (
	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/config"
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/utils"
	"github.com/patrick-jessen/script/lang/jlang/module"
)

type JLang struct {
	workDir string
	Modules []*module.Module
}

func (l *JLang) Compile(path string) []byte {
	f, err := file.Load(path)
	if err != nil {
		utils.ErrLogger.Println(err)
		return nil
	}

	scanner := compiler.NewScanner(l, f)
	tokens := scanner.Scan()
	if config.DebugTokens {
		utils.ErrLogger.Println(compiler.FormatTokens(tokens))
	}

	parser := compiler.NewParser(l, tokens)
	ast := parser.Parse()
	if config.DebugAST {
		utils.ErrLogger.Println(compiler.FormatAST(ast))
	}

	if f.HasErrors() {
		for _, err := range f.Errors {
			utils.ErrLogger.Println(err.Error())
		}
		return nil
	}

	return l.Generate(ast)
}

// func (l *JLang) Compile(path string) []byte {
// 	mainMod, err := module.Load(path, "main")
// 	if err != nil {
// 		utils.ErrLogger.Println(err)
// 		return nil
// 	}

// 	l.workDir = path
// 	l.Modules = []*module.Module{mainMod}

// 	var modMap = map[string]*module.Module{}

// 	// analyze modules
// 	for i := 0; i < len(l.Modules); i++ {
// 		mod := l.Modules[i]
// 		l.analyzeModule(mod)

// 		modMap[mod.Name] = mod
// 	}

// 	// link modules
// 	for _, m := range l.Modules {
// 		for _, i := range m.Imports {
// 			modName := i.Module.Value
// 			symName := i.Symbol.Value

// 			// TODO:
// 			// if strings.HasSuffix(modName, ".dll") {
// 			// 	a.AddExternalSymbol(modName, symName)
// 			// 	continue
// 			// }

// 			mod, ok := modMap[modName]
// 			if !ok {
// 				continue
// 			}

// 			sym, ok := mod.Exports[symName]
// 			if ok {
// 				i.Typ = sym.Type()
// 			} else {
// 				i.Pos().MarkError(fmt.Sprintf(
// 					"module '%v' does not export symbol '%v'",
// 					modName, symName,
// 				))
// 			}
// 		}
// 	}

// 	// check types
// 	a.performTypeCheck(modMap)

// 	if config.DebugAST {
// 		l := log.New(os.Stderr, "", 0)
// 		for k, v := range modMap {
// 			l.Println(color.NewString("AST for module [%v]:", color.Red(k)))
// 			for _, sym := range v.Symbols {
// 				l.Println(ast.FormatAST(sym))
// 			}
// 		}
// 	}

// 	// print errors (if any)
// 	if a.hasErrors() {
// 		a.printErrors()
// 		return fmt.Errorf("compilation aborted due to errors")
// 	}
// 	return nil
// }

// func (l *JLang) analyzeModule(mod *module.Module) {
// 	var symbols = map[string]ast.Declarable{}
// 	var unresolved []*ast.Identifier
// 	var imports []*ast.Identifier

// 	// compile each source file
// 	for _, f := range mod.Files {
// 		scanner := compiler.NewScanner(l, f)
// 		tokens := scanner.Scan()
// 		if config.DebugTokens {
// 			utils.ErrLogger.Println("Tokens for", f.Path)
// 			utils.ErrLogger.Println(compiler.FormatTokens(tokens))
// 		}

// 		parser := compiler.NewParser(l, tokens)
// 		ast := parser.Parse()

// 		// keep track of module-scope symbols
// 		for k, v := range p.Symbols() {
// 			// detect duplicate declarations
// 			if sym, ok := symbols[k]; ok {
// 				v.Pos().MarkError(fmt.Sprintf(
// 					"redeclaration of symbol '%v'. First declared here: (%v)",
// 					k, sym.Pos().Info().Link(),
// 				))
// 			} else {
// 				symbols[k] = v
// 			}
// 		}

// 		// keep track of unresolved symbols (module scope)
// 		unresolved = append(unresolved, p.Unresolved...)

// 		// keep track of imported symbols
// 		imports = append(imports, p.Imports...)

// 		// perform imports
// 		for _, i := range p.ImportedModules {
// 			err := a.importModule(i.Module.Value)
// 			if err != nil {
// 				i.Module.Pos.MarkError(err.Error())
// 			}
// 		}
// 	}

// 	// resolve unresolved references within the module
// 	for _, u := range unresolved {
// 		name := u.Symbol.Value
// 		sym, ok := symbols[name]
// 		if ok {
// 			u.Typ = sym.Type()
// 			u.Obj = sym.Ident().Obj
// 		} else {
// 			u.Pos().MarkError(fmt.Sprintf("unresolved symbol '%v'", name))
// 		}
// 	}

// 	mod.Symbols = symbols
// 	mod.Exports = symbols // TODO: not all symbols should be exported
// 	mod.Imports = imports
// }

// func (l *JLang) importModule(imp string) (err error) {
// 	defer func() {
// 		if e := recover(); e != nil {
// 			err = fmt.Errorf("module '%v' not found", imp)
// 		}
// 	}()

// 	if strings.HasSuffix(imp, ".dll") {
// 		a.AddExternalLib(imp)
// 		return
// 	}

// 	path := path.Join(a.workDir, imp)

// 	// ignore if module is already imported
// 	for _, m := range a.Modules {
// 		if m.Dir == path {
// 			return
// 		}
// 	}
// 	mod := module.Load(path, imp)
// 	a.Modules = append(a.Modules, mod)
// 	return
// }

func (l *JLang) Generate(root ast.Node) []byte {
	return nil
}
