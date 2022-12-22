package jsonlang

import (
	"unicode"

	"github.com/patrick-jessen/script/compiler"
	"github.com/patrick-jessen/script/compiler/ast"
	"github.com/patrick-jessen/script/compiler/config"
	"github.com/patrick-jessen/script/compiler/file"
	"github.com/patrick-jessen/script/compiler/token"
	"github.com/patrick-jessen/script/compiler/utils"
)

type JSONLanguage struct{}

const (
	TokCurlStart = "{"
	TokCurlEnd   = "}"
	TokBracStart = "["
	TokBracEnd   = "]"
	TokColon     = ":"
	TokDot       = "."
	TokComma     = ","
	TokString    = "string"
	TokNumber    = "number"
	TokBool      = "bool"
	TokNull      = "null"
)

func (l *JSONLanguage) Compile(path string) []byte {
	f, err := file.Load(path)
	if err != nil {
		utils.ErrLogger.Println(err)
		return nil
	}

	scanner := compiler.NewScanner(l, f)
	tokens := scanner.Scan()
	if config.DebugTokens {
		utils.ErrLogger.Println(compiler.FormatTokens(tokens))
	}

	parser := compiler.NewParser(l, tokens)
	ast := parser.Parse()
	if config.DebugAST {
		utils.ErrLogger.Println(compiler.FormatAST(ast))
	}

	if f.HasErrors() {
		for _, err := range f.Errors {
			utils.ErrLogger.Println(err.Error())
		}
		return nil
	}

	return l.Generate(ast)
}

func (l *JSONLanguage) Scan(s compiler.LanguageScanner) token.Token {
	if unicode.IsNumber(s.Next()) {
		s.StartCapture()
		for unicode.IsNumber(s.Next()) {
			s.Consume()
		}
		capture1 := s.StopCapture()
		if !s.ConsumeChar('.') {
			return s.TokenVal(TokNumber, capture1)
		}
		s.StartCapture()
		for unicode.IsNumber(s.Next()) {
			s.Consume()
		}
		capture2 := s.StopCapture()
		return s.TokenVal(TokNumber, capture1+"."+capture2)
	}

	if unicode.IsLetter(s.Next()) {
		s.StartCapture()
		for unicode.IsLetter(s.Next()) {
			s.Consume()
		}
		capture := s.StopCapture()
		switch capture {
		case "true", "false":
			return s.TokenVal(TokBool, capture)
		case "null":
			return s.Token(TokNull)
		}
		return s.Error("invalid literal")
	}

	switch s.Next() {
	case ' ', '\t', '\r', '\n':
		s.Consume()
		return s.Skip()
	case '{':
		s.Consume()
		return s.Token(TokCurlStart)
	case '}':
		s.Consume()
		return s.Token(TokCurlEnd)
	case '[':
		s.Consume()
		return s.Token(TokBracStart)
	case ']':
		s.Consume()
		return s.Token(TokBracEnd)
	case ':':
		s.Consume()
		return s.Token(TokColon)
	case ',':
		s.Consume()
		return s.Token(TokComma)
	case '"':
		s.Consume()
		s.StartCapture()
		for !s.NextIs('"') {
			if !s.Consume() {
				break
			}
		}
		capture := s.StopCapture()
		s.Consume()
		return s.TokenVal(TokString, capture)
	}

	s.Consume()
	return s.Error("unknown token")
}

func (l *JSONLanguage) Parse(p compiler.LanguageParser) ast.Node {
	if p.NextIs(TokString) {
		val := p.Next()
		p.Consume()
		return &StringNode{Value: val}
	}
	if p.NextIs(TokNumber) {
		val := p.Next()
		p.Consume()
		return &NumberNode{Value: val}
	}
	if p.NextIs(TokBool) {
		val := p.Next()
		p.Consume()
		return &BoolNode{Value: val}
	}
	if p.NextIs(TokNull) {
		val := p.Next()
		p.Consume()
		return &NullNode{Pos: val.Pos}
	}
	if p.NextIs(TokBracStart) {
		startTok := p.Next()
		var children []ast.Node
		index := 0
		p.Consume()
		for !p.NextIs(TokBracEnd) {
			if p.NextIs(token.EOF) {
				p.Next().Pos.MarkError("unexpected EOF")
				return nil
			}

			if index != 0 {
				if !p.ConsumeType(TokComma) {
					p.Next().Pos.MarkError("expected ,")
				}
			}

			children = append(children, l.Parse(p))
			index++
		}
		if !p.ConsumeType(TokBracEnd) {
			p.Next().Pos.MarkError("expected ]")
		}
		return &ArrayNode{
			Pos:      startTok.Pos,
			Children: children,
		}
	}
	if p.NextIs(TokCurlStart) {
		startTok := p.Next()
		var properties []ast.Node
		isFirstKey := true
		p.Consume()
		for !p.NextIs(TokCurlEnd) {
			if p.NextIs(token.EOF) {
				p.Next().Pos.MarkError("unexpected EOF")
				return nil
			}

			if !isFirstKey {
				if !p.ConsumeType(TokComma) {
					p.Next().Pos.MarkError("expected ,")
				}
			}

			if !p.NextIs(TokString) {
				p.Next().Pos.MarkError("excpected key")
				p.Consume()
				continue
			}

			key := p.Next()
			p.Consume()

			if !p.NextIs(TokColon) {
				p.Next().Pos.MarkError("expected colon")
				p.Consume()
				continue
			}
			p.Consume()

			value := l.Parse(p)
			properties = append(properties, &PropertyNode{
				Name:  key,
				Value: value,
			})
			isFirstKey = false
		}
		if !p.ConsumeType(TokCurlEnd) {
			p.Next().Pos.MarkError("expected }")
		}

		return &ObjectNode{
			Pos:        startTok.Pos,
			Properties: properties,
		}
	}

	return nil
}

func (l *JSONLanguage) Generate(node ast.Node) []byte {
	nodeInfo := node.Info()

	switch nodeInfo.Type {
	case "bool", "number":
		return []byte(nodeInfo.Literal)
	case "string":
		return []byte(`"` + nodeInfo.Literal + `"`)
	case "null":
		return []byte("null")
	case "array":
		out := []byte("[")
		for i, c := range nodeInfo.Children {
			if i != 0 {
				out = append(out, ',')
			}
			out = append(out, l.Generate(c)...)
		}
		out = append(out, ']')
		return out

	case "object":
		out := []byte("{")
		for i, c := range nodeInfo.Children {
			if i != 0 {
				out = append(out, ',')
			}
			out = append(out, l.Generate(c)...)
		}
		out = append(out, '}')
		return out

	case "property":
		var out []byte
		out = append(out, '"')
		out = append(out, []byte(nodeInfo.Name)...)
		out = append(out, '"', ':')
		out = append(out, l.Generate(nodeInfo.Children[0])...)
		return out

	default:
		nodeInfo.Pos.MarkError("unknown node %s", nodeInfo.Type)
		return nil
	}
}
