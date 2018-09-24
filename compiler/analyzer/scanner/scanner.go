// Package scanner is responsible for performing lexical analysis on source code
package scanner

import (
	"fmt"

	"github.com/patrick-jessen/script/config"
	"github.com/patrick-jessen/script/utils/color"
	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

// Scanner is used for tokenizing a source file
type Scanner struct {
	file *file.File // the source file
	iter int        // current index into the source string
	char byte       // current character

	token  token.Token   // current token
	tokens []token.Token // all tokens (if config.DebugTokens == true)
}

// Init initializes the scanner
// Must be called before use
func (s *Scanner) Init(file *file.File) {
	s.file = file
	s.iter = -1 // s.next() increments iter
	s.next()    // read first character
}

// Scan scans the next token
func (s *Scanner) Scan() token.Token {
	var done = false
	for !done && !s.checkEOF() {
		s.token.Value = ""
		s.token.Pos = s.file.NewPos(s.iter)
		done = s.scan()
	}
	if config.DebugTokens {
		s.tokens = append(s.tokens, s.token)
	}
	return s.token
}

// next advances the scanner to the next character in the source string.
// Returns false if EOF is reached, otherwise true.
func (s *Scanner) next() bool {
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

// checkEOF checks if EOF is reached
// returns true if EOF is reached
func (s *Scanner) checkEOF() bool {
	// Scanning while EOF returns the last EOF token
	if s.token.ID == token.EOF {
		s.token.Pos = s.file.NewPos(len(s.file.Source))
		return true
	}
	// note: \n is appended to file
	if s.iter > len(s.file.Source) {
		s.token.Pos = s.file.NewPos(len(s.file.Source))
		s.token.ID = token.EOF
		if config.DebugTokens {
			s.printTokens()
		}
		return true
	}
	return false
}

// printTokens prints the tokens in human readable format
// only valid if config.DebugTokens == true
func (s *Scanner) printTokens() {
	fmt.Println(color.NewString("tokens for [%v]:", color.Red(s.file.Path)))
	for _, t := range s.tokens {
		fmt.Printf("%v\t%v\n", t.Pos.Info().Link(), t)
	}
}
