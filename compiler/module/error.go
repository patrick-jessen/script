package module

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

type SourceError struct {
	Module   Module
	Position int
	Message  string
}

func (se *SourceError) Error() string {
	posInfo := se.Module.PositionInfo(se.Position)

	return fmt.Sprintf(
		"%v\t%v\n%v",
		color.Red("ERROR: "+se.Message),
		posInfo.Link(),
		color.Yellow(posInfo.String()),
	)
}
