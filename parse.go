package main

import (
	"log"
	"strings"
)

func Parse(s string) {

	l := newLexer(s)
	go l.run()

	var ast = newAST()
	var item item
	throw := func() {
		vs := []string{"node:", ast.ptr.typ, "", "type:", string(item.typ)}
		panic(strings.Join(vs, " "))
	}

	for item = range l.items {
		switch item.typ {
		case itemMetaArg:
			switch ast.ptr.typ {
			case "root":
				ast.createChild("article").
					createChild("article.meta").
					createChild("article.meta.args").
					createChild("meta.arg").setValue(item.val).
					back()
			case "article.meta.args":
				ast.createChild("meta.arg").
					setValue(item.val).
					back()
			default:
				throw()
			}
		case itemMetaSep:
			switch ast.ptr.typ {
			case "article.meta.args":
				ast.createSibling("article.meta.kvs")
			default:
				throw()
			}
		case itemMetaKey:
			switch ast.ptr.typ {
			case "article.meta.kvs":
				ast.createChild("meta.kv").
					createChild("meta.key").setValue(item.val)
			default:
				throw()
			}
		case itemMetaValue:
			switch ast.ptr.typ {
			case "meta.key":
				ast.createSibling("meta.value").setValue(item.val).
					back().back()
			default:
				throw()
			}
		case itemBlankLine:
			switch ast.ptr.typ {
			case "article.meta.kvs":
				ast.back().createSibling("article.content").
					createChild("block")
			case "block":
				ast.createChild("blank").back()
			default:
				throw()
			}
		case itemText:
			switch ast.ptr.typ {
			case "block":
				ast.createChild("text").setValue(item.val).back()
			default:
				throw()
			}
		case itemSep:
			switch ast.ptr.typ {
			case "block":
				ast.createChild("sep").back()
			default:
				throw()
			}
		case itemCmd:
			switch ast.ptr.typ {
			case "block":
				ast.createChild("block").
					createChild("block.command").setValue(item.val)
			default:
				throw()
			}
		case itemBlockLeft:
			switch ast.ptr.typ {
			case "block.command":
				ast.back()
			default:
				throw()
			}
		case itemBlockRight:
			switch ast.ptr.typ {
			case "block":
				ast.back()
			default:
				throw()
			}
		case itemEOF:
			ast.back()
			if !ast.is("article.content") {
				log.Println("[Parse]", "EOF at wrong place")
				throw()
			}
		default:
			log.Println("[Parse]", "itemType not support")
			throw()
		}
	}
}
