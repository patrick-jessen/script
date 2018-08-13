package compiler

// func New(dir string) *Compiler {
// 	return &Compiler{
// 		workDir: dir,
// 		modules: []*module.Module{
// 			module.Load(dir, "main"),
// 		},
// 		prog: newProgram(),
// 	}
// }

// func (c *Compiler) generate() {
// 	for _, m := range c.modules {
// 		c.generateModule(m)
// 	}
// }

// func (c *Compiler) generateExpression(n ast.Expression, reg int) (out []ir.Instruction) {

// 	switch exp := n.(type) {
// 	case *ast.String:
// 		l := byte(len(exp.Token.Value))

// 		buf := bytes.Buffer{}
// 		buf.WriteByte(l)
// 		buf.Write(([]byte)(exp.Token.Value))
// 		buf.WriteByte(0)

// 		dPos := c.prog.AddData(buf.Bytes())
// 		out = append(out, &ir.LoadD{Reg: ir.Register(reg), Data: dPos})
// 	case *ast.Integer:
// 		i, _ := strconv.ParseInt(exp.Token.Value, 10, 32)
// 		out = append(out, &ir.LoadI{Reg: ir.Register(reg), Val: int(i)})
// 	case *ast.VariableRef:
// 		if exp.Identifier.Obj.Num < 0 {
// 			out = append(out, &ir.Move{
// 				Dst: ir.Register(reg),
// 				Src: ir.Register(-exp.Identifier.Obj.Num),
// 			})
// 		} else {
// 			fmt.Println("KEK")
// 			panic("Not implemented")
// 		}
// 	default:
// 		fmt.Println("hue", exp, reflect.TypeOf(exp))
// 	}
// 	return
// }

// func (c *Compiler) generateFunction(n *ast.FunctionDecl, modName string) {
// 	fn := &ir.Function{
// 		Name: modName + "." + n.Name(),
// 	}

// 	for i, a := range n.Args {
// 		a.Obj.Num = -i - 1
// 	}

// 	for _, s := range n.Block.Statements {
// 		switch sn := s.(type) {
// 		case *ast.VariableDecl:
// 			sn.Identifier.Obj.Num = fn.NumLocals
// 			fn.NumLocals++

// 			exp := c.generateExpression(sn.Value, 0)

// 			fn.Instructions = append(fn.Instructions, exp...)
// 			fn.Instructions = append(fn.Instructions,
// 				&ir.Set{
// 					Var: ir.Local(sn.Identifier.Obj.Num),
// 					Reg: ir.Register(0),
// 				},
// 			)

// 		case *ast.FunctionCall:
// 			// for i, a := range sn.Args.Args {
// 			// 	exp := c.generateExpression(a, i+1)
// 			// 	fn.Instructions = append(fn.Instructions, exp...)
// 			// }
// 			for i := len(sn.Args.Args) - 1; i >= 0; i-- {
// 				exp := c.generateExpression(sn.Args.Args[i], i+1)
// 				fn.Instructions = append(fn.Instructions, exp...)
// 			}

// 			fn.Instructions = append(fn.Instructions,
// 				&ir.Call{Func: sn.Name()},
// 			)
// 		default:
// 			fmt.Println(reflect.TypeOf(s))
// 		}
// 	}

// 	c.prog.AddFunction(fn)
// }

// func (c *Compiler) generateModule(m *module.Module) {
// 	for _, s := range m.Symbols {
// 		switch n := s.(type) {
// 		case *ast.FunctionDecl:
// 			c.generateFunction(n, m.Name())
// 		}
// 	}
// }
