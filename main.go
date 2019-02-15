package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

func ParseFile(filename string) error {
	t := time.Now()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Println("[Read]", time.Since(t).Nanoseconds(), "ns")
	fmt.Println()

	t = time.Now()
	Parse(string(content))
	fmt.Println("[Parse]", time.Since(t).Nanoseconds(), "ns")
	return nil
}

func throw(ast *node, item item) {
	vs := []string{"node:", ast.ptr.typ, "", "type:", string(item.typ)}
	panic(strings.Join(vs, " "))
}

func Parse(s string) {
	l := newLexer(s)
	go l.run()

	ast := newAST()
	for item := range l.items {
		fmt.Println(item.typ, item.val)
		switch item.typ {
		case itemMetaArg:
			switch ast.ptr.typ {
			case "root":
				ast.createChild("article").
					createChild("article.meta").
					createChild("article.meta.args").
					createChild("meta.arg").
					setValue(item.val).
					back()
			case "article.meta.args":
				ast.createChild("meta.arg").
					setValue(item.val).
					back()
			default:
				throw(ast, item)
			}
		case itemMetaSep:
			switch ast.ptr.typ {
			case "article.meta.args":
				ast.createSibling("article.meta.kvs")
			default:
				throw(ast, item)
			}
		case itemMetaKey:
			switch ast.ptr.typ {
			case "article.meta.kvs":
				ast.createChild("meta.kv").
					createChild("meta.key").
					setValue(item.val)
			default:
				throw(ast, item)
			}
		case itemMetaValue:
			switch ast.ptr.typ {
			case "meta.key":
				ast.createSibling("meta.value").
					setValue(item.val).
					back().
					back()
			default:
				throw(ast, item)
			}
		case itemText:
		}
	}
}

func main() {
	_ = ParseFile("test.rtcl")
}
