package main

import (
	"fmt"
	"os"

	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/linker"
	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/vm"
)

type command int

const compVer = "0.1"
const (
	cmdRun command = iota
	cmdDebug
	cmdBuild
)

var (
	dir = "./src" // TODO: change to ./
	cmd = cmdRun
)

func printHelp() {
	fmt.Println(
		"Script ver. " + compVer + "\n" +
			"Usage:\n" +
			"    script <command> [srcDir]\n" +
			"Commands:\n" +
			"    run\n" +
			"        runs the application\n" +
			"    debug\n" +
			"        debugs the application\n" +
			"    build\n" +
			"        builds the application to an executable\n",
	)
}

func handleArgs() bool {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "run":
			cmd = cmdRun
		case "build":
			cmd = cmdBuild
		case "debug":
			cmd = cmdDebug
		default:
			fmt.Println(color.Red("Invalid command\n"))
			printHelp()
			return false
		}
	} else {
		printHelp()
		return false
	}

	if len(os.Args) > 2 {
		dir = os.Args[2]
		if _, err := os.Stat(dir); err != nil {
			fmt.Println(color.Red("Invalid path\n"))
			printHelp()
			return false
		}
	}
	return true
}

func main() {
	if !handleArgs() {
		return
	}

	comp := compiler.New(dir)
	prog := comp.Run()

	switch cmd {
	case cmdRun:
		vm.Run(prog, false)
	case cmdDebug:
		vm.Run(prog, true)
	case cmdBuild:
		linker.Run(prog)
	}
}
