package utils

import (
	"fmt"
	"strings"

	"github.com/patrick-jessen/script/compiler/config"
)

type color int

const (
	white color = iota
	red
	green
	yellow
	blue
)

func White(str interface{}) String {
	return String{
		string: fmt.Sprintf("%v", str),
		color:  white,
	}
}
func Red(str interface{}) String {
	color := red
	if config.NoColor {
		color = white
	}
	return String{
		string: fmt.Sprintf("%v", str),
		color:  color,
	}
}
func Green(str interface{}) String {
	color := green
	if config.NoColor {
		color = white
	}
	return String{
		string: fmt.Sprintf("%v", str),
		color:  color,
	}
}
func Yellow(str interface{}) String {
	color := yellow
	if config.NoColor {
		color = white
	}
	return String{
		string: fmt.Sprintf("%v", str),
		color:  color,
	}
}
func Blue(str interface{}) String {
	color := blue
	if config.NoColor {
		color = white
	}
	return String{
		string: fmt.Sprintf("%v", str),
		color:  color,
	}
}

type String struct {
	subStrings []String
	string     string
	color      color
}

func NewString(fmt string, a ...String) String {
	split := strings.Split(fmt, "%v")
	colStr := String{}

	for i, s := range split {
		colStr.subStrings = append(colStr.subStrings, White(s))
		if len(a) > i {
			colStr.subStrings = append(colStr.subStrings, a[i])
		}
	}

	return colStr
}

func (cs *String) Add(s String) {
	cs.subStrings = append(cs.subStrings, s)
}

func (cs *String) Length() int {
	l := len(cs.string)
	for _, s := range cs.subStrings {
		l += s.Length()
	}
	return l
}

func (cs String) String() (out string) {
	switch cs.color {
	case white:
		out += cs.string
	case red:
		out += fmt.Sprintf("\033[1;31m%v\033[0m", cs.string)
	case green:
		out += fmt.Sprintf("\033[1;32m%v\033[0m", cs.string)
	case yellow:
		out += fmt.Sprintf("\033[1;33m%v\033[0m", cs.string)
	case blue:
		out += fmt.Sprintf("\033[1;34m%v\033[0m", cs.string)
	}

	for _, s := range cs.subStrings {
		out += s.String()
	}
	return
}