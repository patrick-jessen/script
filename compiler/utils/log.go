package utils

import (
	"log"
	"os"
)

var ErrLogger = log.New(os.Stderr, "", 0)
