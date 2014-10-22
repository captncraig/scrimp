package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int
type Pos int
type Token struct {
	Type TokenType
	Pos  Pos
	Val  string
}

const (
	TokError TokenType = iota
	TokEOF
	TokIdentifier
	TokInteger
	TokNamespace
	TokInclude
	TokConst
	TokTypedef
	TokEnum
	TokStruct
	TokException
	TokService
	TokColon
	TokComma
	TokLCurly
	TokRCurly
	TokSemicolon
	TokLParen
	TokRParen
	TokDocText
	TokRequired
	TokOptional
	TokOneway
	TokThrows
	TokVoid
	TokLAngle
	TokRAngle
)

var idTokens = map[string]TokenType{
	"namespace": TokNamespace,
	"include":   TokInclude,
	"const":     TokConst,
	"typedef":   TokTypedef,
	"enum":      TokEnum,
	"struct":    TokStruct,
	"exception": TokException,
	"service":   TokService,
	"required":  TokRequired,
	"optional":  TokOptional,
	"oneway":    TokOneway,
	"throws":    TokThrows,
}

const eof = -1

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input   string     // the string being scanned
	pos     Pos        // current position in the input
	start   Pos        // start position of this item
	width   Pos        // width of last rune read from input
	lastPos Pos        // position of most recent item returned by nextItem
	tokens  chan Token // channel of scanned items
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// lineNumber reports which line we're on, based on the position of
// the previous item returned by nextItem. Doing it this way
// means we don't have to worry about peek double counting.
func (l *Lexer) lineNumber() int {
	return 1 + strings.Count(l.input[:l.lastPos], "\n")
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- Token{TokError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next item from the input.
func (l *Lexer) nextToken() Token {
	item := <-l.tokens
	l.lastPos = item.Pos
	return item
}

func Lex(input string) <-chan Token {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
	}
	go l.run()
	return l.tokens
}

func (l *Lexer) run() {
	fmt.Println("yyy")
	s := startState
	for s != nil {
		s = s(l)
	}

}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func startState(l *Lexer) stateFn {
	r := l.next()
	if r == '_' || unicode.IsLetter(r) {
		return lexIdentifier
	}
	l.emit(TokEOF)
	return nil
}

func lexIdentifier(l *Lexer) stateFn {
	for {
		r := l.next()
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' {
			//absorb
		} else {
			l.backup()
			emitIdentifier(l.input[l.start:l.pos], l)
			return startState
		}
	}
}

func emitIdentifier(word string, l *Lexer) {
	typ, ok := idTokens[word]
	if ok {
		l.emit(typ)
		return
	}
	l.emit(TokIdentifier)

}
