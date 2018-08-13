package file

import (
	"fmt"

	"github.com/patrick-jessen/script/utils/color"
)

// Error describes an error within a file
type Error struct {
	Position Pos    // the position of the error
	Message  string // the error message
}

// Error implements the error interface
func (se *Error) Error() string {
	// get info regarding the error position
	posInfo := se.Position.Info()

	return fmt.Sprintf(
		"%v\t%v\n%v",
		color.Red("ERROR: "+se.Message),
		posInfo.Link(),
		color.Yellow(posInfo),
	)
}
