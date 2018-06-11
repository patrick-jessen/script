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
		0x48, 0xc7, 0xc0, 0x2a, 0x00, 0x00, 0x00, 0x48, 0x83, 0xc0, 0x0a, 0xc3,
	}

	p := pe.New(0x1000, 0x200)
	p.SetCode(code)

	p.WriteFile("./out.exe")

	fmt.Println("time:", time.Since(start))
	fmt.Println("file written: ./out.exe")
}
