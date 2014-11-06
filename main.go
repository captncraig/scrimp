package main

import (
	"fmt"
	"github.com/captncraig/scrimp/model"
	"github.com/captncraig/scrimp/parse"
	"log"
	"os"
	"text/template"
)

func main() {
	program := `
	namespace csharp foobar
	namespace java foo2bar namespace * blah
	namespace * blah2
	include "abc"
	include 'a"bc'
	const list<map<Foo,bar>> foo="abc"
	typedef int Foo
	typedef list<map<Foo,bar>> complexListType
	`
	doc := parse.Parse(program)
	fmt.Println(doc)
	for _, c := range doc.Consts {
		fmt.Println(c)
	}
	for _, t := range doc.Typedefs {
		fmt.Println(t)
	}

}

func testGenerate() {
	temp, err := template.ParseFiles("html/program.tpl", "html/doc.tpl", "html/toc.tpl")
	for _, t := range temp.Templates() {
		fmt.Println(t.Name())
	}
	if err != nil {
		log.Fatalln(err)
	}
	w, _ := os.Create("out.html")

	p := model.TProgram{}
	p.Name = "tutorial"
	p.SetDoc(`The first thing to know about are types. The available types in Thrift are:

bool        Boolean, one byte
byte        Signed byte
i16         Signed 16-bit integer
i32         Signed 32-bit integer
i64         Signed 64-bit integer
double      64-bit floating point value
string      String
binary      Blob (byte array)
map<t1,t2>  Map from one type to another
list<t1>    Ordered list of one type
set<t1>     Set of unique elements of one type
Did you also notice that Thrift supports C style comments?`)
	s := model.TService{}
	s.Name = "Calculator"
	f := model.TFunction{}
	f.Name = "Add"
	s.AddFunc(&f)
	p.AddService(&s)

	temp.Execute(w, p)
}
