package main

import (
	"fmt"
	"github.com/captncraig/scrimp/parse"
)

func main() {
	program := `
	namespace csharp foobar 
	namespace java foo2bar namespace * blah
	namespace * blah2
	`
	fmt.Println(parse.Parse(program))
}
