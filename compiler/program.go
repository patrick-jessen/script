package compiler

import (
	"crypto/md5"
	"fmt"

	"github.com/patrick-jessen/script/compiler/ir"
	"github.com/patrick-jessen/script/utils/color"
)

type Program struct {
	Functions  map[string]*ir.Function
	Data       []byte
	DataCache  map[[16]byte]int
	SharedLibs map[string][]string
}

func newProgram() *Program {
	return &Program{
		Functions:  make(map[string]*ir.Function),
		DataCache:  make(map[[16]byte]int),
		SharedLibs: make(map[string][]string),
	}
}

func (p *Program) AddExternalLib(libName string) {
	if _, ok := p.SharedLibs[libName]; !ok {
		p.SharedLibs[libName] = []string{}
	}
}
func (p *Program) AddExternalSymbol(libName string, symName string) {
	p.SharedLibs[libName] = append(p.SharedLibs[libName], symName)
}

func (p *Program) AddFunction(fn *ir.Function) {
	p.Functions[fn.Name] = fn
}
func (p *Program) AddData(d []byte) (out int) {
	hash := md5.Sum(d)
	if p, ok := p.DataCache[hash]; ok {
		return p
	}
	out = len(p.Data)
	p.Data = append(p.Data, d...)
	p.DataCache[hash] = out
	return
}

func (p *Program) String() (out string) {
	out += fmt.Sprintf("Functions:  %v\t[", len(p.Functions))
	first := true
	for _, f := range p.Functions {
		if !first {
			out += ", "
		} else {
			first = false
		}
		out += color.Red(f.Name)
	}
	out += "]\n"
	out += fmt.Sprintf("DataItems:  %v\t(%v bytes)\n", len(p.DataCache), len(p.Data))
	out += fmt.Sprintf("SharedLibs: %v\t[", len(p.SharedLibs))

	first = true
	for k := range p.SharedLibs {
		if !first {
			out += ", "
		} else {
			first = false
		}
		out += color.Yellow(k)
	}
	out += "]"

	return
}
