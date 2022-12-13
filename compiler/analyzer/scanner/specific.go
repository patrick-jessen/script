package scanner

import (
	"unicode"

	"github.com/patrick-jessen/script/utils/token"
)

// Scan scans the next token
// returns false to skip token
// NOTE: consider creating interface for language spec
func (s *Scanner) scan() bool {
	if isLetter(s.char) {
		i := s.scanIdentifer()
		k := keywordLookUp(i)
		if k != token.Invalid {
			// Token is a keyword
			s.token.ID = k
		} else {
			s.token.ID = token.Identifier
			s.token.Value = i
		}
		return true
	}
	if isDigit(s.char) {
		v := s.scanNumber()
		if s.char == '.' {
			s.next()
			v2 := s.scanNumber()
			s.token.ID = token.Float
			s.token.Value = v + "." + v2
		} else {
			s.token.ID = token.Integer
			s.token.Value = v
		}
		return true
	}

	switch s.char {
	case ' ':
		fallthrough
	case '\t':
		s.next()
		return false
	case '(':
		s.token.ID = token.ParentStart
	case ')':
		s.token.ID = token.ParentEnd
	case '{':
		s.token.ID = token.CurlStart
	case '}':
		s.token.ID = token.CurlEnd
	case '=':
		s.token.ID = token.Equal
	case '+':
		s.token.ID = token.Plus
	case '-':
		s.token.ID = token.Minus
	case '*':
		s.token.ID = token.Asterisk
	case '/':
		s.token.ID = token.Slash
		s.next()
		if s.char == '/' || s.char == '*' {
			s.scanComment()
			return false
		}
		return true
	case ',':
		s.token.ID = token.Comma
	case '.':
		s.token.ID = token.Dot
	case '"':
		str := s.scanString()
		s.token.ID = token.String
		s.token.Value = str
	case '\r':
		s.next()
		fallthrough
	case '\n':
		s.token.ID = token.NewLine
	default:
		s.file.NewPos(s.iter).MarkError("unexpected token")
		s.next()
		return false
	}
	s.next()
	return true
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
			s.file.NewPos(s.iter).MarkError("expected \"")
			break
		}
	}
	return string(s.file.Source[start+1 : s.iter])
}

func (s *Scanner) scanIdentifer() string {
	start := s.iter
	for isLetter(s.char) {
		s.next()
	}
	return string(s.file.Source[start:s.iter])
}

func (s *Scanner) scanNumber() string {
	start := s.iter
	for isDigit(s.char) {
		s.next()
	}
	return string(s.file.Source[start:s.iter])
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

func isLetter(b rune) bool {
	return unicode.IsLetter(b)
}

func isDigit(b rune) bool {
	return unicode.IsNumber(b)
}
