package main

import (
	"encoding/json"
	"fmt"

	"github.com/patrick-jessen/script/src/parser"
)

func main() {
	ast, err := parser.New("src.j").Parse()
	if err != nil {
		fmt.Println(err)
		return
	}

	out, _ := json.MarshalIndent(ast, "", "  ")
	fmt.Println(string(out))
}
