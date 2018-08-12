package token

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/utils/color"
)

// ID is a unique identifier for types of tokens.
type ID int

// Token is an actual token which has been read from the input string.
type Token struct {
	Pos   file.Pos // Position of token
	ID    ID       // ID of the token
	Value string   // The captured value of the token (not all types of tokens have values)
}

// String returns the name of the token ID
func (i ID) String() string {
	return names[i]
}

// Name returns the name of the token ID
func (t Token) Name() string {
	return t.ID.String()
}

// String returns a pretty formatting of the token
func (t Token) String() string {
	return color.NewString("%v %v", color.Green(t.ID), color.Yellow(t.Value)).String()
}

var names = map[ID]string{
	Invalid:     "INVALID",
	EOF:         "EOF",
	Var:         "var",
	Func:        "func",
	Import:      "import",
	Return:      "return",
	NewLine:     "newline",
	Identifier:  "identifier",
	Equal:       "=",
	Plus:        "+",
	Minus:       "-",
	Asterisk:    "*",
	Slash:       "/",
	ParentStart: "(",
	ParentEnd:   ")",
	CurlStart:   "{",
	CurlEnd:     "}",
	Comma:       ",",
	Dot:         ".",
	Float:       "float",
	Integer:     "integer",
	String:      "string",
}

const (
	Invalid ID = iota
	EOF
	Var
	Func
	Import
	Return
	NewLine
	Identifier
	Equal
	Plus
	Minus
	Asterisk
	Slash
	ParentStart
	ParentEnd
	CurlStart
	CurlEnd
	Comma
	Dot
	Float
	Integer
	String
)
