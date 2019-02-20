package rtcl

import (
	"io/ioutil"
	"log"
	"strings"
)

func Parse(s string) *node {
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
					createChild("article.meta.Args").
					createChild("meta.arg").setValue(item.val).
					back()
			case "article.meta.Args":
				ast.createChild("meta.arg").
					setValue(item.val).
					back()
			default:
				throw()
			}
		case itemMetaSep:
			switch ast.ptr.typ {
			case "article.meta.Args":
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
					createChild("block").
					createChild("block.command").setValue("_wrapper").
					back()
			case "block":
				ast.createChild("blank").back()
			default:
				throw()
			}
		case itemText:
			switch ast.ptr.typ {
			case "block":
				child := ast.ptr.child
				if child != nil {
					for ; child.sibling != nil; child = child.sibling {
					}
				}

				if child != nil && child.typ == "paragraph" {
					ast.ptr = child
				} else {
					ast.createChild("paragraph")
				}
				ast.createChild("text").setValue(item.val).back()
				ast.back()
			default:
				throw()
			}
		//case itemSep:
		//	switch ast.ptr.typ {
		//	case "block":
		//		ast.createChild("sep").back()
		//	default:
		//		throw()
		//	}
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
		case itemMetaItem:
			switch ast.ptr.typ {
			case "block":
				ast.createChild("meta.item").setValue(item.val)
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
	ast.backToRoot()
	return ast
}

func ParseFile(filename string) (ast *node, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	ast = Parse(string(content))
	return
}
