package jlang

import (
	"unicode"

	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/utils/token"

	"github.com/patrick-jessen/script/lang/jlang/tokens"
)

func (l *JLang) Scan(s compiler.LanguageScanner) token.Token {
	if isLetter(s.Next()) {
		return l.scanIdentifer(s)
	}
	if isDigit(s.Next()) {
		return l.scanNumber(s)
	}

	switch s.Next() {
	case ' ':
		s.Consume()
		return s.Skip()
	case '\t':
		s.Consume()
		return s.Skip()
	case '(':
		s.Consume()
		return s.Token(tokens.ParentStart)
	case ')':
		s.Consume()
		return s.Token(tokens.ParentEnd)
	case '{':
		s.Consume()
		return s.Token(tokens.CurlStart)
	case '}':
		s.Consume()
		return s.Token(tokens.CurlEnd)
	case '=':
		s.Consume()
		return s.Token(tokens.Equal)
	case '+':
		s.Consume()
		return s.Token(tokens.Plus)
	case '-':
		s.Consume()
		return s.Token(tokens.Minus)
	case '*':
		s.Consume()
		return s.Token(tokens.Asterisk)
	case '/':
		s.Consume()
		if s.Next() == '/' || s.Next() == '*' {
			return l.scanComment(s)
		}
		return s.Token(tokens.Slash)
	case ',':
		s.Consume()
		return s.Token(tokens.Comma)
	case '.':
		s.Consume()
		return s.Token(tokens.Dot)
	case '"':
		return l.scanString(s)
	case '\r':
		s.Consume()
		return s.Skip()
	case '\n':
		s.Consume()
		return s.Skip()
	default:
		s.Consume()
		return s.Error("unexpected token")
	}
}

func (l *JLang) scanComment(s compiler.LanguageScanner) token.Token {
	// first slash is already consumed
	if s.Next() == '/' {
		for s.Next() != '\n' {
			if !s.Consume() {
				break
			}
		}
		return s.Skip()
	}

	depth := 1

	for depth != 0 {
		if !s.Consume() {
			break
		}
		if s.Next() == '/' {
			s.Consume()
			if s.Next() == '*' {
				s.Consume()
				depth++
			}
		} else if s.Next() == '*' {
			s.Consume()
			if s.Next() == '/' {
				s.Consume()
				depth--
			}
		}
	}
	return s.Skip()
}

func (l *JLang) scanString(s compiler.LanguageScanner) token.Token {
	s.Consume()
	s.StartCapture()
	for s.Next() != '"' {
		if !s.Consume() {
			return s.Error("expected \"")
		}
	}
	capture := s.StopCapture()
	s.Consume()
	return s.TokenVal(tokens.String, capture)
}

func (l *JLang) scanIdentifer(s compiler.LanguageScanner) token.Token {
	s.StartCapture()
	for isLetter(s.Next()) {
		s.Consume()
	}
	capture := s.StopCapture()
	var typ string

	switch capture {
	case "var":
		typ = tokens.Var
	case "func":
		typ = tokens.Func
	case "import":
		typ = tokens.Import
	case "return":
		typ = tokens.Return
	default:
		return s.TokenVal(tokens.Identifier, capture)
	}

	return s.Token(typ)
}

func (l *JLang) scanNumber(s compiler.LanguageScanner) token.Token {
	s.StartCapture()
	for isDigit(s.Next()) {
		s.Consume()
	}
	capture1 := s.StopCapture()

	if s.Next() != '.' {
		return s.TokenVal(tokens.Integer, capture1)
	}

	s.Consume()
	s.StartCapture()
	for isDigit(s.Next()) {
		s.Consume()
	}
	capture2 := s.StopCapture()

	return s.TokenVal(tokens.Float, capture1+"."+capture2)
}

func isLetter(b rune) bool {
	return unicode.IsLetter(b)
}

func isDigit(b rune) bool {
	return unicode.IsNumber(b)
}
