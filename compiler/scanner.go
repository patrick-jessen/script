package compiler

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

// Scanner is used for tokenizing a source file
type Scanner struct {
	lang       ScannerImpl
	file       *file.File  // the source file
	iter       int         // current index into the source string
	captures   []int       // start indices of captures
	char       rune        // current character
	token      token.Token // current token
	tokenStart int
}

// NewScanner creates a new scanner
func NewScanner(lang ScannerImpl, file *file.File) *Scanner {
	s := &Scanner{file: file, lang: lang}
	s.Reset()
	return s
}

// Reset moves the iterator to the begging of the file
func (s *Scanner) Reset() {
	s.iter = -1
	s.token = token.Token{}
	s.Consume()
}

// Scan scans the entire file
func (s *Scanner) Scan() []token.Token {
	var out []token.Token
	t := s.NextToken()
	for {
		out = append(out, t)
		if t.Type == token.EOF {
			break
		}
		t = s.NextToken()
	}
	return out
}

// NextToken scans the next token
func (s *Scanner) NextToken() token.Token {
	for !s.checkEOF() {
		s.tokenStart = s.iter
		s.token = s.lang.Scan(s)
		s.token.Pos = s.file.NewPos(s.tokenStart)
		s.captures = nil

		if s.token.Type != token.Skip {
			break
		}
	}
	return s.token
}

// Consume advances the scanner to the next character in the source string.
// Returns false if EOF is reached, otherwise true.
func (s *Scanner) Consume() bool {
	s.iter++
	// check for EOF
	if s.iter >= len(s.file.Source) {
		// insert newline at the end of file to make parsing easier
		s.char = '\n'
		return false
	}
	// advance to next character
	s.char = s.file.Source[s.iter]
	// mark newlines - this is done to help the file resolve token positions
	if s.char == '\n' {
		s.file.MarkLine(s.iter)
	}
	return true
}

// ConsumeChar attempts to consume a specific character
func (s *Scanner) ConsumeChar(char rune) bool {
	if s.NextIs(char) {
		return s.Consume()
	}
	return false
}

// Next returns the next character in the source
func (s *Scanner) Next() rune {
	return s.char
}

// NextIs returns whether the next character in the source is of a certain value
func (s *Scanner) NextIs(char rune) bool {
	return s.char == char
}

// Token creates a new token
func (s *Scanner) Token(typ string) token.Token {
	return token.Token{
		Type: typ,
	}
}

// TokenVal creates a new token with value
func (s *Scanner) TokenVal(typ string, value string) token.Token {
	return token.Token{
		Type:  typ,
		Value: value,
	}
}

// Skip creates a new skip token
func (s *Scanner) Skip() token.Token {
	return token.Token{Type: token.Skip}
}

// Error creates a new error token
func (s *Scanner) Error(format string, args ...interface{}) token.Token {
	s.file.NewPos(s.tokenStart).MarkError(format, args...)
	return s.Skip()
}

// StartCapture starts capturing characters from current position
func (s *Scanner) StartCapture() {
	s.captures = append(s.captures, s.iter)
}

// StopCapture returns the captured string
func (s *Scanner) StopCapture() string {
	idx := len(s.captures) - 1
	start := s.captures[idx]
	s.captures = s.captures[:idx]
	return string(s.file.Source[start:s.iter])
}

// checkEOF checks if EOF is reached
// returns true if EOF is reached
func (s *Scanner) checkEOF() bool {
	// Scanning while EOF returns the last EOF token
	if s.token.Type == token.EOF {
		s.token.Pos = s.file.NewPos(len(s.file.Source))
		return true
	}
	// note: \n is appended to file
	if s.iter > len(s.file.Source) {
		s.token.Pos = s.file.NewPos(len(s.file.Source))
		s.token.Type = token.EOF
		return true
	}
	return false
}
