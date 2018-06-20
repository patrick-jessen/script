package file

import (
	"fmt"

	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/utils/color"
)

// fileError describes an error within a file
type fileError struct {
	File     *File     // the file in which the error occured
	Position token.Pos // the position of the error
	Message  string    // the error message
}

// Error makes fileError implement error interface
func (se *fileError) Error() string {
	// get info regarding the error position
	posInfo := se.File.PosInfo(se.Position)

	return fmt.Sprintf(
		"%v\t%v\n%v",
		color.Red("ERROR: "+se.Message),
		posInfo.Link(),
		color.Yellow(posInfo),
	)
}
