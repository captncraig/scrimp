package parse

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
		case TokEOF:
			return doc
		default:
			panic("Unexpected token!!!")
		}
	}
	return doc
}

func (p *parser) parseNamespace(m map[string]string) {
	lang := p.require(TokIdentifier, TokStar)
	namespace := p.require(TokIdentifier)
	m[lang.Val] = namespace.Val
}
