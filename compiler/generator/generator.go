package generator

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/patrick-jessen/script/compiler/analyzer"

	"github.com/patrick-jessen/script/compiler/analyzer/ast"
	"github.com/patrick-jessen/script/compiler/generator/ir"
	"github.com/patrick-jessen/script/utils/module"
)

type Generator struct {
	modules []*module.Module
	prog    *Program
}

func New(a *analyzer.Analyzer) *Generator {
	return &Generator{
		modules: a.Modules,
		prog:    newProgram(a.SharedLibs),
	}
}

func (g *Generator) Run() *Program {
	for _, m := range g.modules {
		g.generateModule(m)
	}
	return g.prog
}

func (g *Generator) generateExpression(n ast.Expression, reg int) (out []ir.Instruction) {

	switch exp := n.(type) {
	case *ast.String:
		l := byte(len(exp.Token.Value))

		buf := bytes.Buffer{}
		buf.WriteByte(l)
		buf.Write(([]byte)(exp.Token.Value))
		buf.WriteByte(0)

		dPos := g.prog.AddData(buf.Bytes())
		out = append(out, &ir.LoadD{Reg: ir.Register(reg), Data: dPos})
	case *ast.Integer:
		i, _ := strconv.ParseInt(exp.Token.Value, 10, 32)
		out = append(out, &ir.LoadI{Reg: ir.Register(reg), Val: int(i)})
	case *ast.VariableRef:
		if exp.Identifier.Obj.Num < 0 {
			out = append(out, &ir.Move{
				Dst: ir.Register(reg),
				Src: ir.Register(-exp.Identifier.Obj.Num),
			})
		} else {
			fmt.Println("KEK")
			panic("Not implemented")
		}
	default:
		fmt.Println("hue", exp, reflect.TypeOf(exp))
	}
	return
}

func (g *Generator) generateFunction(n *ast.FunctionDecl, modName string) {
	fn := &ir.Function{
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

			exp := g.generateExpression(sn.Value, 0)

			fn.Instructions = append(fn.Instructions, exp...)
			fn.Instructions = append(fn.Instructions,
				&ir.Set{
					Var: ir.Local(sn.Identifier.Obj.Num),
					Reg: ir.Register(0),
				},
			)

		case *ast.FunctionCall:
			// for i, a := range sn.Args.Args {
			// 	exp := c.generateExpression(a, i+1)
			// 	fn.Instructions = append(fn.Instructions, exp...)
			// }
			for i := len(sn.Args.Args) - 1; i >= 0; i-- {
				exp := g.generateExpression(sn.Args.Args[i], i+1)
				fn.Instructions = append(fn.Instructions, exp...)
			}

			fn.Instructions = append(fn.Instructions,
				&ir.Call{Func: sn.Name()},
			)
		default:
			fmt.Println(reflect.TypeOf(s))
		}
	}

	g.prog.AddFunction(fn)
}

func (g *Generator) generateModule(m *module.Module) {
	for _, s := range m.Symbols {
		switch n := s.(type) {
		case *ast.FunctionDecl:
			g.generateFunction(n, m.Name())
		}
	}
}
