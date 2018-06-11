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
		//0xB8, 0x2A, 0x00, 0x00, 0x00, // MOV eax, 42
		//0xC3, // RET

		0x48, 0xc7, 0xc0, 0x2a, 0x00, 0x00, 0x00, 0x48, 0x83, 0xc0, 0x0a, 0xc3,
	}

	p := pe.New(0x1000, 0x200)

	textSec := p.NewSection(".text")
	textSec.SetData(code)

	// p := pe.PE{
	// 	LinkerVersion:    [2]byte{0, 1},
	// 	FileAlignment:    512,
	// 	SectionAlignment: 512,

	// 	Sections: []pe.Section{
	// 		pe.Section{
	// 			Name:          [8]byte{'.', 't', 'e', 'x', 't', 0, 0, 0},
	// 			SizeOfRawData: int32(len(code)),
	// 		},
	// 		pe.Section{
	// 			Name: [8]byte{'.', 'd', 'a', 't', 'a', 0, 0, 0},
	// 		},
	// 	},
	// }
	p.WriteFile("./out.exe")

	fmt.Println("time:", time.Since(start))
	fmt.Println("file written: ./out.exe")
}
