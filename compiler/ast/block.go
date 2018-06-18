package ast

type Block struct {
	Statements []Node
}

func (b *Block) String() (out string) {
	for _, s := range b.Statements {
		out += s.String() + "\n"
	}
	return
}
