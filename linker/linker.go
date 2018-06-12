package linker

import (
	"fmt"
	"time"

	pe "github.com/patrick-jessen/script/linker/PE"
)

func Run() {
	fmt.Println("LINKING ==================================")
	start := time.Now()

	code := []byte{
		// return 42 + 10
		// 0x48, 0xc7, 0xc0, 0x2a, 0x00, 0x00, 0x00, 0x48, 0x83, 0xc0, 0x0a, 0xc3,

		0x48, 0xC7, 0xC1, 0x2A, 0x00, 0x00, 0x00, 0x48, 0x83, 0xC1, 0x0A, 0xFF, 0x14, 0x25, 0x78, 0x20, 0x40, 0x00,
	}

	p := pe.New(0x1000, 0x200)
	p.SetCode(code)
	p.Import("ExitProcess", "kernel32.dll")
	p.Import("MessageBoxA", "user32.DLL")
	p.WriteFile("./out.exe")

	fmt.Println("time:", time.Since(start))
	fmt.Println("file written: ./out.exe")
}
