package compiler

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/module"
	"github.com/patrick-jessen/script/compiler/parser"
	"github.com/patrick-jessen/script/vm"
)

type SharedLib struct {
	Lib     string
	Symbols []string
}

type Compiler struct {
	workDir string
	modules []*module.Module

	sharedLibs []*SharedLib

	pm *vm.ProgManager
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

	if strings.HasSuffix(imp, ".dll") {
		for _, l := range c.sharedLibs {
			if l.Lib == imp {
				return
			}
		}
		c.sharedLibs = append(c.sharedLibs, &SharedLib{Lib: imp})
		return
	}

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
	defer func() {
		// recover()
		c.printErrors()
	}()

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

			if strings.HasSuffix(modName, ".dll") {
				var sl *SharedLib
				for _, l := range c.sharedLibs {
					if l.Lib == modName {
						sl = l
						break
					}
				}
				sl.Symbols = append(sl.Symbols, symName)
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

	// for _, sl := range c.sharedLibs {
	// 	fmt.Println(*sl)
	// }

	c.generate()
}

func (c *Compiler) generate() {
	c.pm = vm.NewPM()

	for _, m := range c.modules {
		c.generateModule(m)
	}

	c.pm.AddFunction(&vm.Function{
		Name: "fmt.print",
		Instructions: []vm.Instruction{
			&vm.CallGo{
				Func: func(vm *vm.VM) int {
					first := vm.GoString(vm.Get(1))
					second := vm.GoString(vm.Get(2))

					fmt.Println(first, second)
					return 0
				},
			},
		},
	})

	c.pm.Run()
}

func (c *Compiler) generateExpression(n ast.Expression, reg int) (out []vm.Instruction) {

	switch exp := n.(type) {
	case *ast.String:
		l := byte(len(exp.Token.Value))
		dPos := c.pm.AddBytes([]byte{l})
		c.pm.AddBytes(([]byte)(exp.Token.Value))
		out = append(out, &vm.LoadC{Dst: reg, Val: dPos})
	case *ast.Integer:
		i, _ := strconv.ParseInt(exp.Token.Value, 10, 32)
		out = append(out, &vm.LoadC{Dst: reg, Val: int(i)})
	case *ast.VariableRef:
		if exp.Identifier.Obj.Num < 0 {
			out = append(out, &vm.Mov{
				Dst: reg,
				Src: -exp.Identifier.Obj.Num,
			})
		} else {
			out = append(out,
				&vm.Get{Local: exp.Identifier.Obj.Num, Reg: reg})
		}
	default:
		fmt.Println("hue", exp, reflect.TypeOf(exp))
	}
	return
}

func (c *Compiler) generateFunction(n *ast.FunctionDecl, modName string) {
	fn := &vm.Function{
		Name: modName + "." + n.Name(),
	}

	for i, a := range n.Args {
		a.Obj.Num = -i - 1
	}

	for _, s := range n.Block.Statements {
		switch sn := s.(type) {
		case *ast.VariableDecl:
			sn.Identifier.Obj.Num = fn.NumLocals
			fn.NumLocals++

			exp := c.generateExpression(sn.Value, 0)

			fn.Instructions = append(fn.Instructions, exp...)
			fn.Instructions = append(fn.Instructions,
				&vm.Set{Local: sn.Identifier.Obj.Num, Reg: 0},
			)

		case *ast.FunctionCall:
			for i, a := range sn.Args.Args {
				exp := c.generateExpression(a, i+1)
				fn.Instructions = append(fn.Instructions, exp...)
			}

			fn.Instructions = append(fn.Instructions,
				&vm.Call{Func: sn.Name()})
		default:
			fmt.Println(reflect.TypeOf(s))
		}
	}

	c.pm.AddFunction(fn)
	fmt.Println(fn)
}

func (c *Compiler) generateModule(m *module.Module) {
	for _, s := range m.Symbols {
		switch n := s.(type) {
		case *ast.FunctionDecl:
			c.generateFunction(n, m.Name())
		}
	}

	// c.pm.Print()
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
			err := c.importModule(i.Module.Value)
			if err != nil {
				mod.Error(i.Module.Pos, err.Error())
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
			mod.Error(u.Pos(), fmt.Sprintf("unresolved symbol '%v'", name))
		}
	}

	mod.Symbols = symbols
	mod.Exports = symbols // TODO: not all symbols should be exported
	mod.Imports = imports
}
