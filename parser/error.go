package parser

import (
	"fmt"
	"strings"
)

type parserError struct {
	Message  string
	Parser   *Parser
	Position int
}

func (e *parserError) Error() string {
	src := string(e.Parser.src)
	lines := strings.Split(src, "\n")
	var actualLine string

	iter := 0
	line := 0

	for iter <= e.Position {
		l := lines[line]
		len := len(l)

		if iter+len > e.Position {
			actualLine = l
			break
		} else {
			iter += len
			line++
		}
	}

	lineNumb := fmt.Sprintf("%v:", line+1)
	pos := fmt.Sprintf("%v\t%v^", strings.Repeat(" ", len(lineNumb)), strings.Repeat(" ", e.Position))

	return fmt.Sprintf("ERROR: %v\n%v\t%v\n%v", e.Message, lineNumb, actualLine, pos)
}
