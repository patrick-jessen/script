package token

// Pos is the position of a token in a file.
type Pos int

// ID is a unique identifier for types of tokens.
type ID int

// Token is an actual token which has been read from the input string.
type Token struct {
	Pos   Pos    // Position of token
	ID    ID     // ID of the token
	Value string // The captured value of the token (not all types of tokens have values)
}

func (i ID) String() string {
	return names[i]
}
func (t Token) Name() string {
	return t.ID.String()
}

var names = map[ID]string{
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
	Float:       "float",
	Integer:     "integer",
	String:      "string",
}

const (
	Var ID = iota
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
	Float
	Integer
	String
)
