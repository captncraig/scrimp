package parse

import (
	"fmt"
	"strings"
)

type parser struct {
	lex    <-chan Token
	peeked *Token
}

func Parse(input string) *Document {
	l := Lex(input)
	parser := parser{lex: l}
	return parser.parseDocument()
}

func (p *parser) nextToken() Token {
	if p.peeked != nil {
		tok := *p.peeked
		p.peeked = nil
		return tok
	}
	tok := <-p.lex
	return tok
}

func (p *parser) peek() Token {
	if p.peeked != nil {
		return *p.peeked
	}
	tok := <-p.lex
	p.peeked = &tok
	return tok
}

func (p *parser) require(typ ...TokenType) Token {
	tok := p.nextToken()
	for _, tt := range typ {
		if tok.Type == tt {
			return tok
		}
	}
	panic("unexpected token")
}

func (p *parser) parseDocument() *Document {
	doc := NewDocument()
	for {
		tok := p.nextToken()
		switch tok.Type {
		case TokNamespace:
			p.parseNamespace(doc.Namespaces)
		case TokInclude:
			p.parseInclude(doc)
		case TokConst:
			doc.AddConst(p.parseConst())
		case TokTypedef:
			doc.AddTypedef(p.parseTypeDef())
		case TokEOF:
			return doc
		default:
			panic("Unexpected token!!!")
		}
	}
	return doc
}

func (p *parser) parseTypeDef() *Typedef {
	return &Typedef{
		FieldType: p.parseFieldType(),
		Name:      p.require(TokIdentifier).Val,
	}
}

func (p *parser) parseConst() *Constant {
	c := Constant{}
	c.FieldType = p.parseFieldType()
	c.Name = p.require(TokIdentifier).Val
	p.require(TokEqual)
	c.Value = p.parseConstValue()
	fmt.Println(c.Value)
	return &c
}

func (p *parser) parseConstValue() string {
	tok := p.require(TokNumConst, TokStringConst)
	return tok.Val
}

func (p *parser) parseFieldType() string {
	tok := p.require(TokIdentifier, TokSet, TokMap, TokList)
	if tok.Type == TokIdentifier {
		return tok.Val
	}
	if tok.Type == TokSet {
		p.require(TokLAngle)
		inner := p.parseFieldType()
		p.require(TokRAngle)
		return "set<" + inner + ">"
	}
	if tok.Type == TokList {
		p.require(TokLAngle)
		inner := p.parseFieldType()
		p.require(TokRAngle)
		return "list<" + inner + ">"
	}
	if tok.Type == TokMap {
		p.require(TokLAngle)
		key := p.parseFieldType()
		p.require(TokComma)
		val := p.parseFieldType()
		p.require(TokRAngle)
		return "map<" + key + "," + val + ">"
	}
	panic("cant get here")
}

func (p *parser) parseInclude(d *Document) {
	q := p.require(TokStringConst)
	d.Includes = append(d.Includes, strings.Trim(q.Val, "\"'"))
}

func (p *parser) parseNamespace(m map[string]string) {
	lang := p.require(TokIdentifier, TokStar)
	namespace := p.require(TokIdentifier)
	m[lang.Val] = namespace.Val
}
