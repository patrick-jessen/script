// Package scanner is responsible for tokenizing source files
package scanner

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

// Scanner is used for tokenizing a source file
type Scanner struct {
	file *file.File // the source file
	iter int        // current index into the source string
	char byte       // current character
}

// Init initializes the scanner
// Must be called before use
func (s *Scanner) Init(file *file.File) {
	s.file = file
	s.iter = 0
	if len(file.Source) > 0 {
		// point to first character
		s.char = file.Source[0]
	} else {
		// insert newline at the end of file
		s.char = '\n'
	}
}

// next advances the scanner to the next character in the source string.
// Returns false if EOF is reached, otherwise true.
func (s *Scanner) next() bool {
	s.iter++

	if s.iter < len(s.file.Source) {
		// advance to next character
		s.char = s.file.Source[s.iter]

		// mark newlines - this is done to help the file resolve token positions
		if s.char == '\n' {
			s.file.MarkLine(s.iter)
		}
		return true
	}

	// EOF is reached
	s.char = '\n' // insert newline at the end of file
	return false
}

// Scan scans the next token.
func (s *Scanner) Scan() (tok token.Token) {
startScan:
	// If EOF is reached, return 'EOF' token
	if s.iter > len(s.file.Source) {
		return token.Token{
			ID:  token.EOF,
			Pos: s.file.Pos(len(s.file.Source)),
		}
	}

	tok.Pos = s.file.Pos(s.iter)

	if isLetter(s.char) {
		i := s.scanIdentifer()
		k := keywordLookUp(i)
		if k != token.Invalid {
			// Token is a keyword
			tok.ID = k
		} else {
			tok.ID = token.Identifier
			tok.Value = i
		}
		return
	}
	if isDigit(s.char) {
		v := s.scanNumber()
		if s.char == '.' {
			s.next()
			v2 := s.scanNumber()
			tok.ID = token.Float
			tok.Value = v + "." + v2
		} else {
			tok.ID = token.Integer
			tok.Value = v
		}
		return
	}

	switch s.char {
	case ' ':
		fallthrough
	case '\t':
		s.next()
		goto startScan
	case '(':
		tok.ID = token.ParentStart
	case ')':
		tok.ID = token.ParentEnd
	case '{':
		tok.ID = token.CurlStart
	case '}':
		tok.ID = token.CurlEnd
	case '=':
		tok.ID = token.Equal
	case '+':
		tok.ID = token.Plus
	case '-':
		tok.ID = token.Minus
	case '*':
		tok.ID = token.Asterisk
	case '/':
		tok.ID = token.Slash
		s.next()
		if s.char == '/' || s.char == '*' {
			s.scanComment()
			goto startScan
		}
		return
	case ',':
		tok.ID = token.Comma
	case '.':
		tok.ID = token.Dot
	case '"':
		str := s.scanString()
		tok.ID = token.String
		tok.Value = str
	case '\r':
		s.next()
		fallthrough
	case '\n':
		tok.ID = token.NewLine
	default:
		s.file.Error(token.Pos(s.iter), "unexpected token")
		s.next()
		goto startScan
	}

	s.next()
	return
}

func (s *Scanner) scanComment() {
	// first slash is already consumed
	if s.char == '/' {
		for s.char != '\n' {
			if !s.next() {
				break
			}
		}
	} else {
		counter := 1

		for counter != 0 {
			if !s.next() {
				break
			}
			if s.char == '/' {
				s.next()
				if s.char == '*' {
					s.next()
					counter++
				}
			} else if s.char == '*' {
				s.next()
				if s.char == '/' {
					s.next()
					counter--
				}
			}
		}
	}
}

func (s *Scanner) scanString() string {
	start := s.iter
	s.next()
	for s.char != '"' {
		if !s.next() {
			s.file.Error(s.file.Pos(s.iter),
				"expected \"")
			break
		}
	}
	return s.file.Source[start+1 : s.iter]
}

func (s *Scanner) scanIdentifer() string {
	start := s.iter
	for isLetter(s.char) {
		s.next()
	}
	return s.file.Source[start:s.iter]
}

func (s *Scanner) scanNumber() string {
	start := s.iter
	for isDigit(s.char) {
		s.next()
	}
	return s.file.Source[start:s.iter]
}

func keywordLookUp(str string) token.ID {
	if str == "var" {
		return token.Var
	}
	if str == "func" {
		return token.Func
	}
	if str == "import" {
		return token.Import
	}
	if str == "return" {
		return token.Return
	}
	return token.Invalid
}

func isLetter(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z'
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}
