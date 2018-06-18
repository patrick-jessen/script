package scanner

import (
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
)

type Scanner struct {
	iter int
	file *file.File
	char byte
}

func (s *Scanner) Init(file *file.File) {
	s.file = file
	s.char = file.Source[0]
}

func (s *Scanner) next() {
	s.iter++
	if s.iter < len(s.file.Source) {
		s.char = s.file.Source[s.iter]
	}
}

func (s *Scanner) Scan() (tok token.Token) {
startScan:
	if s.iter == len(s.file.Source) {
		return token.Token{
			ID:  token.EOF,
			Pos: token.Pos(len(s.file.Source)) | s.file.PosMask,
		}
	}

	tok.Pos = token.Pos(s.iter) | s.file.PosMask

	if isLetter(s.char) {
		i := s.scanIdentifer()
		k := keywordLookUp(i)
		if k != token.Invalid {
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
	case ',':
		tok.ID = token.Comma
	case '"':
		str := s.scanString()
		tok.ID = token.String
		tok.Value = str
	case '\r':
		s.next()
		fallthrough
	case '\n':
		tok.ID = token.NewLine
		s.file.MarkLine(s.iter)
	default:
		s.file.Error(token.Pos(s.iter), "unexpected token")
		s.next()
		goto startScan
	}

	s.next()
	return
}

func (s *Scanner) scanString() string {
	start := s.iter
	s.next()
	for s.char != '"' {
		s.next()
	}
	s.next()
	return s.file.Source[start+1 : s.iter-1]
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
