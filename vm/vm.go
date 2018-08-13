package vm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/patrick-jessen/script/compiler/generator"
	"github.com/patrick-jessen/script/compiler/generator/ir"
	"github.com/patrick-jessen/script/utils/color"
)

func Run(prog *generator.Program, debug bool) int {
	vm := newVM(prog, debug)
	vm.Call("main.main")
	return vm.regs[0]
}

type vm struct {
	debug  bool
	stack  []int
	regs   [8]int
	prog   *generator.Program
	extern map[string]uintptr
}

func newVM(prog *generator.Program, debug bool) *vm {
	vm := &vm{
		prog:   prog,
		extern: make(map[string]uintptr),
		debug:  debug,
	}

	for libName, syms := range prog.SharedLibs {
		lib, err := syscall.LoadLibrary(libName)
		if err != nil {
			panic(err)
		}

		for _, s := range syms {
			proc, err := syscall.GetProcAddress(lib, s)
			if err != nil {
				panic(err)
			}
			vm.extern[libName+"."+s] = proc
		}
	}

	return vm
}

func (vm *vm) printDebug(fn *ir.Function, inst int) {

	var instLines []color.String
	var regsLines []color.String
	var stackLines []color.String
	var dataLines []color.String

	instLines = append(instLines, color.White("== Instructions ============================================"))
	instLines = append(instLines, color.NewString("%v %v", color.Blue("func"), color.Red(fn.Name)))

	for idx, i := range fn.Instructions {
		str := color.NewString("")

		if idx == inst {
			str.Add(color.White("> "))
		} else {
			str.Add(color.White("  "))
		}
		str.Add(i.ColorString())
		instLines = append(instLines, str)
	}

	regsLines = append(regsLines, color.White("== Registers ================="))

	for i := 0; i < 8; i++ {
		regsLines = append(regsLines, color.NewString(
			"reg%v  %v", color.Red(i), color.White(vm.regs[i]),
		))
	}

	stackLines = append(stackLines, color.White("== Stack ====================="))
	for i := 0; i < fn.NumLocals; i++ {
		stackLines = append(stackLines, color.NewString(
			"loc%v  %v", color.Blue(i), color.White(vm.stack[i]),
		))
	}

	dataLines = append(dataLines, color.White("== Data ==================="))
	numDataLines := int(float32(len(vm.prog.Data))/8 + 0.5)
	for i := 0; i < numDataLines; i++ {
		line := color.String{}
		for ii := 0; ii < 8; ii++ {
			line.Add(color.White(""))
		}

		dataLines = append(dataLines, color.White("test"))
	}

	max := len(instLines)
	if len(regsLines) > max {
		max = len(regsLines)
	}
	if len(stackLines) > max {
		max = len(stackLines)
	}
	if len(dataLines) > max {
		max = len(dataLines)
	}

	for i := 0; i < max; i++ {
		column := 0
		fmt.Print("| ")

		if len(instLines) > i {
			fmt.Print(instLines[i])
			column = instLines[i].Length()
		}
		fmt.Print(strings.Repeat(" ", 60-column))
		fmt.Print(" | ")

		if len(regsLines) > i {
			fmt.Print(regsLines[i])
			column = regsLines[i].Length()
		} else {
			column = 0
		}

		fmt.Print(strings.Repeat(" ", 30-column))
		fmt.Print(" | ")

		if len(stackLines) > i {
			fmt.Print(stackLines[i])
			column = stackLines[i].Length()
		} else {
			column = 0
		}

		fmt.Print(strings.Repeat(" ", 30-column))
		fmt.Print(" | ")

		if len(dataLines) > i {
			fmt.Print(dataLines[i])
			column = dataLines[i].Length()
		} else {
			column = 0
		}
		fmt.Print(strings.Repeat(" ", 30-column))
		fmt.Println(" |")
	}
}

func (vm *vm) Call(fnName string) {
	fn, ok := vm.prog.Functions[fnName]
	if !ok {
		// perhaps it is external
		vm.callC(fnName)
		return
	}

	bp := len(vm.stack)
	newStack := make([]int, bp+fn.NumLocals)
	copy(newStack, vm.stack)
	vm.stack = newStack

	for i := 0; i < len(fn.Instructions); i++ {
		if vm.debug {
			fmt.Println(strings.Repeat("\n", 10))
			vm.printDebug(fn, i)
			reader := bufio.NewReader(os.Stdin)
			reader.ReadBytes('\n')
		}
		fn.Instructions[i].Execute(vm)
	}

	vm.stack = vm.stack[:bp]
}

func (vm *vm) Reg(reg ir.Register) int {
	return vm.regs[int(reg)]
}
func (vm *vm) SetReg(reg ir.Register, val int) {
	vm.regs[int(reg)] = val
}
func (vm *vm) Stack(loc ir.Local) int {
	return vm.stack[len(vm.stack)-int(loc)-1]
}
func (vm *vm) SetStack(loc ir.Local, val int) {
	vm.stack[len(vm.stack)-int(loc)-1] = val
}
func (vm *vm) Data(idx int) int {
	return int(uintptr(unsafe.Pointer(&vm.prog.Data[0]))) + idx
}

func (vm *vm) callC(fnName string) {
	fn, ok := vm.extern[fnName]
	if !ok {
		panic(fmt.Sprintf("missing function '%v'", fnName))
	}

	ret, _, err := syscall.Syscall6(fn, 4,
		uintptr(vm.regs[1]), uintptr(vm.regs[2]),
		uintptr(vm.regs[3]), uintptr(vm.regs[4]),
		0, 0)

	if err != 0 {
		fmt.Printf("cannot call external function '%v'\n", fnName)
	}
	vm.regs[0] = int(ret)
}
