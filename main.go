package main

import (
	"fmt"
	"github.com/captncraig/scrimp/lexer"
)

func main() {
	lex := lexer.Lex("include")
	for {
		tok := <-lex
		fmt.Println(tok)
		if tok.Type == lexer.TokEOF {
			return
		}

	}
}
