package nodes

type FunctionCallArgs struct {
	Args []Expression
}

// func (f *FunctionCallArgs) Pos() file.Pos {
// 	return f.Args[0].Pos()
// }
// func (f *FunctionCallArgs) TypeCheck() {
// 	for _, a := range f.Args {
// 		a.TypeCheck()
// 	}
// }