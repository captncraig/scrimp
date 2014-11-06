package parse

import (
	"fmt"
	"strconv"
	"strings"
)

type parser struct {
	lex     <-chan Token
	peeked  *Token
	docText string
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
	tok := p.grabToken()
	return tok
}

func (p *parser) grabToken() Token {
	tok := <-p.lex
	for tok.Type == TokDocText {
		fmt.Println(tok.Val)
		tok.Val = tok.Val[3 : len(tok.Val)-2]
		tok.Val = strings.Trim(tok.Val, " \t\r\n")
		p.docText += tok.Val
		tok = <-p.lex
	}
	return tok
}

func (p *parser) takeIf(t TokenType) bool {
	if p.peek().Type == t {
		p.nextToken()
		return true
	}
	return false
}

func (p *parser) peek() Token {
	if p.peeked != nil {
		return *p.peeked
	}
	tok := p.grabToken()
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
		if p.takeIf(TokNamespace) {
			p.acceptProgramDoctext(doc)
			p.parseNamespace(doc.Namespaces)
		} else if p.takeIf(TokInclude) {
			p.acceptProgramDoctext(doc)
			p.parseInclude(doc)
		} else {
			break
		}
	}
	for {
		tok := p.nextToken()
		switch tok.Type {
		case TokConst:
			doc.AddConst(p.parseConst())
		case TokTypedef:
			doc.AddTypedef(p.parseTypeDef())
		case TokEnum:
			doc.AddEnum(p.parseEnum())
		case TokStruct:
			doc.AddStruct(p.parseStruct())
		case TokException:
			doc.AddXception(p.parseStruct())
		case TokService:
			doc.AddService(p.parseService())
		case TokEOF:
			return doc
		default:
			fmt.Println(tok.Type)
			panic("Unexpected token!!!")
		}
		p.clearDocText() //so docs don't leak from declaration to declaration if misplaced
	}
	return doc
}

func (p *parser) parseTypeDef() *Typedef {
	return &Typedef{
		DocText:   p.clearDocText(),
		FieldType: p.parseFieldType(),
		Name:      p.require(TokIdentifier).Val,
	}
}

func (p *parser) acceptProgramDoctext(d *Document) {
	if p.docText != "" {
		d.DocText = p.clearDocText()
	}
}

func (p *parser) clearDocText() string {
	s := p.docText
	p.docText = ""
	return s
}

func (p *parser) parseService() *Service {
	s := Service{Name: p.require(TokIdentifier).Val, Functions: []*Function{}}
	s.DocText = p.clearDocText()
	if p.takeIf(TokExtends) {
		s.Extends = p.require(TokIdentifier).Val
	}
	p.require(TokLCurly)
	for !p.takeIf(TokRCurly) {
		s.AddFunction(p.parseFunction())
		if p.peek().Type == TokComma || p.peek().Type == TokSemicolon {
			p.nextToken()
		}
	}
	return &s
}

func (p *parser) parseFunction() *Function {
	f := Function{}
	f.DocText = p.clearDocText()
	if p.takeIf(TokOneway) {
		f.Oneway = true
	}
	if p.takeIf(TokVoid) {
		f.ReturnType = "void"
	} else {
		f.ReturnType = p.parseFieldType()
	}
	f.Name = p.require(TokIdentifier).Val
	p.require(TokLParen)
	f.Fields = p.parseFieldList()
	p.require(TokRParen)
	if p.takeIf(TokThrows) {
		p.require(TokLParen)
		f.Throws = p.parseFieldList()
		p.require(TokRParen)
	}
	p.clearDocText()
	return &f
}

func (p *parser) parseStruct() *Struct {
	s := Struct{Name: p.require(TokIdentifier).Val}
	s.DocText = p.clearDocText()
	p.require(TokLCurly)
	s.Fields = p.parseFieldList()
	p.require(TokRCurly)
	return &s
}

func (p *parser) parseFieldList() []*Field {
	nextIdx := 1
	fields := []*Field{}
	for {
		f := p.parseField(nextIdx)
		if f == nil {
			break
		}
		fields = append(fields, f)
		nextIdx = f.Index + 1
	}
	return fields
}

func (p *parser) parseField(nextIdx int) *Field {
	tok := p.peek()
	if tok.Type == TokRCurly || tok.Type == TokRParen {
		return nil
	}
	if tok.Type == TokNumConst {
		p.nextToken()
		p.require(TokColon)
		i, err := strconv.ParseInt(tok.Val, 10, 32)
		if err != nil || int(i) < nextIdx {
			panic("Invalid integer for enum member index")
		}
		nextIdx = int(i)
	}
	f := Field{Index: nextIdx}
	f.DocText = p.clearDocText()
	if p.takeIf(TokRequired) {
		f.Required = true
	} else {
		p.takeIf(TokOptional)
	}
	f.FieldType = p.parseFieldType()
	f.Name = p.require(TokIdentifier).Val
	if p.takeIf(TokEqual) {
		f.DefaultValue = p.parseConstValue()
	}
	if p.peek().Type == TokSemicolon || p.peek().Type == TokComma {
		p.nextToken()
	}
	return &f
}

func (p *parser) parseEnum() *Enum {
	en := Enum{Name: p.require(TokIdentifier).Val, Members: map[string]int{}}
	en.DocText = p.clearDocText()
	p.require(TokLCurly)
	curIdx := 1
	for {
		if p.peek().Type == TokIdentifier {
			name := p.nextToken().Val
			//Does it have an index specified?
			if p.takeIf(TokEqual) {
				num := p.require(TokNumConst).Val
				i, err := strconv.ParseInt(num, 10, 32)
				if err != nil || int(i) < curIdx {
					panic("Invalid integer for enum member index")
				}
				curIdx = int(i)
			}
			en.Members[name] = curIdx
			curIdx++
			if p.peek().Type == TokComma || p.peek().Type == TokSemicolon {
				p.nextToken()
			}
		} else {
			p.require(TokRCurly)
			break
		}
	}
	return &en
}

func (p *parser) parseConst() *Constant {
	c := Constant{}
	c.DocText = p.clearDocText()
	c.FieldType = p.parseFieldType()
	c.Name = p.require(TokIdentifier).Val
	p.require(TokEqual)
	c.Value = p.parseConstValue()
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
