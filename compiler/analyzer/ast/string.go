package ast

import (
	"strings"

	"github.com/patrick-jessen/script/utils/file"
	"github.com/patrick-jessen/script/utils/token"
)

type String struct {
	Token token.Token
}

func (s *String) String(level int) (out string) {
	out = s.Token.Pos.Info().Link()
	out += strings.Repeat("  ", level)
	out += s.Token.String()
	out += "\n"
	return
}

func (s *String) Pos() file.Pos {
	return s.Token.Pos
}

func (s *String) Type() Type {
	return Type{
		IsResolved: true,
		Return:     "string",
	}
}
func (*String) TypeCheck() {
}
