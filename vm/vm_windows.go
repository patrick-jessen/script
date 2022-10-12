package vm

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
