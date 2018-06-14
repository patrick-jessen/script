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
