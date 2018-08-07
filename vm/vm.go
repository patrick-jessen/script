package vm

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/patrick-jessen/script/compiler/ir"

	"github.com/patrick-jessen/script/compiler"
)

func Run(prog *compiler.Program, debug bool) int {
	vm := newVM(prog)
	vm.Call("main.main")
	return vm.regs[0]
}

type vm struct {
	stack  []int
	regs   [8]int
	prog   *compiler.Program
	extern map[string]uintptr
}

func newVM(prog *compiler.Program) *vm {
	vm := &vm{
		prog:   prog,
		extern: make(map[string]uintptr),
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

	fmt.Println(vm.regs)

	ret, _, err := syscall.Syscall6(fn, 4,
		uintptr(vm.regs[1]), uintptr(vm.regs[0]),
		uintptr(vm.regs[0]), uintptr(vm.regs[4]),
		0, 0)

	if err != 0 {
		fmt.Printf("cannot call external function '%v'\n", fnName)
	}
	vm.regs[0] = int(ret)
}
