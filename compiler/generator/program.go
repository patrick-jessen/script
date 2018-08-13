package generator

import (
	"crypto/md5"

	"github.com/patrick-jessen/script/compiler/generator/ir"
)

type Program struct {
	Functions  map[string]*ir.Function
	Data       []byte
	DataCache  map[[16]byte]int
	SharedLibs map[string][]string
}

func newProgram(sharedLibs map[string][]string) *Program {
	return &Program{
		Functions:  make(map[string]*ir.Function),
		DataCache:  make(map[[16]byte]int),
		SharedLibs: sharedLibs,
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
