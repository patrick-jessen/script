package module

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type sourceError struct {
	Module   *Module
	Position token.Pos
	Message  string
}

func (se *sourceError) Error() string {
	posInfo := se.Module.PosInfo(se.Position)

	return fmt.Sprintf(
		"%v\t%v\n%v",
		color.Red("ERROR: "+se.Message),
		posInfo.Link(),
		color.Yellow(posInfo.String()),
	)
}
