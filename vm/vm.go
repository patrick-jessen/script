package vm

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

type ProgManager struct {
	data []byte
	fns  map[string]*Function
}

func NewPM() *ProgManager {
	return &ProgManager{
		fns: make(map[string]*Function),
	}
}
func (pm *ProgManager) AddBytes(d []byte) int {
	l := len(pm.data)
	pm.data = append(pm.data, d...)
	return l
}
func (pm *ProgManager) AddFunction(f *Function) {
	pm.fns[f.Name] = f
}
func (pm *ProgManager) Print() {
	for _, fn := range pm.fns {
		fmt.Println(fn)
	}
}

type Function struct {
	Name         string
	NumLocals    int
	Instructions []Instruction
}

func (f *Function) String() (out string) {
	out = fmt.Sprintf("%v %v",
		color.Blue("func"), color.Red(f.Name))

	if f.NumLocals > 0 {
		out += fmt.Sprintf(" (%v locals)", f.NumLocals)
	}
	out += "\n"

	for i := 0; i < len(f.Instructions); i++ {
		out += "  " + f.Instructions[i].String() + "\n"
	}
	return
}

func (f *Function) Execute(vm *VM) {
	bp := len(vm.stack)
	newStack := make([]int, bp+f.NumLocals)
	copy(newStack, vm.stack)
	vm.stack = newStack

	for i := 0; i < len(f.Instructions); i++ {
		f.Instructions[i].Execute(vm)
	}
	vm.stack = vm.stack[:bp]
}

var data = []byte{
	5, 'h', 'e', 'l', 'l', 'o', 0,
	5, 'w', 'o', 'r', 'l', 'd', 0,
}

var functions = map[string]*Function{}

var print = &Function{
	Name:      "print",
	NumLocals: 0,
	Instructions: []Instruction{
		&CallGo{
			Func: func(vm *VM) int {
				l := int(vm.data[vm.regs[1]])
				s := vm.regs[1] + 1
				e := s + l
				first := string(vm.data[s:e])

				l = int(vm.data[vm.regs[2]])
				s = vm.regs[2] + 1
				e = s + l
				second := string(vm.data[s:e])

				fmt.Println(first, second)
				return 0
			},
		},
	},
}

var main_main = &Function{
	Name:      "main",
	NumLocals: 1,
	Instructions: []Instruction{
		&LoadC{Dst: 7, Val: 0},
		&Set{Local: 0, Reg: 7},

		&Get{Local: 0, Reg: 1},
		&LoadC{Dst: 2, Val: 7},
		&Call{Func: "print"},
	},
}

func (pm *ProgManager) Run() {
	vm := New(pm.data)
	vm.pm = pm
	pm.fns["main.main"].Execute(vm)
	fmt.Println("returned", vm.regs[0])
}

type VM struct {
	stack []int
	regs  [8]int
	data  []byte
	pm    *ProgManager
}

func New(data []byte) *VM {
	return &VM{data: data}
}

type Push struct {
	Reg int
}

func (p *Push) Execute(vm *VM) {
	vm.stack = append(vm.stack, vm.regs[p.Reg])
}
func (vm *VM) Get(reg int) int {
	return vm.regs[reg]
}
func (vm *VM) GoString(data int) string {
	l := int(vm.data[data])
	s := data + 1
	e := s + l
	return string(vm.data[s:e])
}

type Pop struct {
	Reg int
}

func (p *Pop) Execute(vm *VM) {
	l := len(vm.stack) - 1
	vm.regs[p.Reg] = vm.stack[l]
	vm.stack = vm.stack[:l]
}

type Instruction interface {
	Execute(vm *VM)
	String() string
}

type Mov struct {
	Dst int
	Src int
}

func (m *Mov) Execute(vm *VM) {
	vm.regs[m.Dst] = vm.regs[m.Src]
}
func (m *Mov) String() string {
	return fmt.Sprintf("%v reg%v reg%v",
		color.Yellow("Mov"), color.Red(m.Dst), color.Red(m.Src))
}

type Add struct {
	Dst int
	Src int
}

func (a *Add) Execute(vm *VM) {
	vm.regs[a.Dst] = vm.regs[a.Dst] + vm.regs[a.Src]
}
func (*Add) String() string {
	return "NA"
}

type LoadC struct {
	Dst int
	Val int
}

func (c *LoadC) Execute(vm *VM) {
	vm.regs[c.Dst] = c.Val
}
func (c *LoadC) String() string {
	return fmt.Sprintf("%v\treg%v %v",
		color.Yellow("LoadC"),
		color.Red(c.Dst),
		color.Green(c.Val),
	)
}

type CallGo struct {
	Func func(vm *VM) int
}

func (c *CallGo) Execute(vm *VM) {
	vm.regs[0] = c.Func(vm)
}
func (*CallGo) String() string {
	return fmt.Sprintf("%v [native]", color.Yellow("CallGo"))
}

type Set struct {
	Local int
	Reg   int
}

func (s *Set) Execute(vm *VM) {
	vm.stack[len(vm.stack)-1-s.Local] = vm.regs[s.Reg]
}
func (s *Set) String() string {
	return fmt.Sprintf("%v\tloc%v reg%v",
		color.Yellow("Set"),
		color.Blue(s.Local),
		color.Red(s.Reg),
	)
}

type Get struct {
	Local int
	Reg   int
}

func (g *Get) Execute(vm *VM) {
	vm.regs[g.Reg] = vm.stack[len(vm.stack)-1-g.Local]
}
func (g *Get) String() string {
	return fmt.Sprintf("%v\treg%v loc%v",
		color.Yellow("Get"),
		color.Red(g.Reg),
		color.Blue(g.Local),
	)
}

type Call struct {
	Func string
}

func (c *Call) Execute(vm *VM) {
	vm.pm.fns[c.Func].Execute(vm)
}

func (c *Call) String() string {
	return fmt.Sprintf("%v\t%v",
		color.Yellow("Call"), color.Red(c.Func))
}
