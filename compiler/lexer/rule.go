package lexer

import (
	"regexp"
)

// Rule is a rule for matching tokens.
type Rule struct {
	TokenID TokenID
	Name    string
	regexp  *regexp.Regexp
	Omit    bool
}

// NewRule creates a new token rule.
func NewRule(tokenID TokenID, name string, regSrc string) Rule {
	return Rule{
		TokenID: tokenID,
		Name:    name,
		regexp:  regexp.MustCompile("^" + regSrc),
	}
}

// NewOmitRule creates a new token rule which will be omitted from the
// token stream.
func NewOmitRule(name string, regSrc string) Rule {
	return Rule{
		Name:   name,
		regexp: regexp.MustCompile("^" + regSrc),
		Omit:   true,
	}
}
