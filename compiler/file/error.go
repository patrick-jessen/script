package file

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

type fileError struct {
	File     *File
	Position token.Pos
	Message  string
}

func (se *fileError) Error() string {
	posInfo := se.File.PosInfo(se.Position)

	return fmt.Sprintf(
		"%v\t%v\n%v",
		color.Red("ERROR: "+se.Message),
		posInfo.Link(),
		color.Yellow(posInfo.String()),
	)
}
