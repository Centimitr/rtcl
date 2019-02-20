package rtcl

import (
	"io/ioutil"
	"log"
	"rtcl/ast"
	"strings"
)

func Parse(s string) *node {
	l := newLexer(s)
	go l.run()

	var cur = newASTCursor(newAST())
	//cur.onCreate = func() {
	//	fmt.Println(strings.Repeat("    ", cur.ptr.depth()-1) + cur.ptr.typ)
	//}

	var item item
	throw := func() {
		vs := []string{"node:", cur.typ, "", "type:", string(item.typ)}
		panic(strings.Join(vs, " "))
	}

	for item = range l.items {
		switch item.typ {
		case itemMetaArg:
			switch cur.typ {
			case ast.Root:
				cur.createChild(ast.Article).
					createChild(ast.ArticleMeta).
					createChild(ast.ArticleMetaArgs).
					createChild(ast.MetaArg).setValue(item.val).
					back()
			case ast.ArticleMetaArgs:
				cur.createChild(ast.MetaArg).
					setValue(item.val).
					back()
			default:
				throw()
			}
		case itemMetaSep:
			switch cur.typ {
			case ast.ArticleMetaArgs:
				cur.createSibling(ast.ArticleMetaKVs)
			default:
				throw()
			}
		case itemMetaKey:
			switch cur.typ {
			case ast.ArticleMetaKVs:
				cur.createChild(ast.MetaKV).
					createChild(ast.MetaKey).setValue(item.val)
			default:
				throw()
			}
		case itemMetaValue:
			switch cur.typ {
			case ast.MetaKey:
				cur.createSibling(ast.MetaValue).setValue(item.val).
					back().back()
			default:
				throw()
			}
		case itemBlankLine:
			switch cur.typ {
			case ast.ArticleMetaKVs:
				cur.back().createSibling(ast.ArticleContent).
					createChild(ast.Block).
					createChild(ast.BlockCommand).setValue(ast.CommandWrapper).
					back()
			case ast.Block:
				cur.createChild(ast.Blank).back()
			default:
				throw()
			}
		case itemText:
			switch cur.typ {
			case ast.Block:
				clone := cur.clone()
				if clone.gotoLastChild() && clone.typ == ast.Paragraph {
					cur.set(clone.ptr)
				} else {
					cur.createChild(ast.Paragraph)
				}
				cur.createChild(ast.Text).setValue(item.val).back()
				cur.back()
			default:
				throw()
			}
		//case itemSep:
		//	switch cur.ptr.typ {
		//	case "block":
		//		cur.createChild("sep").back()
		//	default:
		//		throw()
		//	}
		case itemCmd:
			switch cur.typ {
			case ast.Block:
				cur.createChild(ast.Block).
					createChild(ast.BlockCommand).setValue(item.val)
			default:
				throw()
			}
		case itemBlockLeft:
			switch cur.typ {
			case ast.BlockCommand:
				cur.back()
			default:
				throw()
			}
		case itemBlockRight:
			switch cur.typ {
			case ast.Block:
				cur.back()
			default:
				throw()
			}
		case itemMetaItem:
			switch cur.typ {
			case ast.Block:
				cur.createChild(ast.MetaItem).setValue(item.val)
				cur.back()
			default:
				throw()
			}
		case itemEOF:
			cur.back()
			if cur.typ != ast.ArticleContent {
				log.Println("[Parse]", "EOF at wrong place")
				throw()
			}
		default:
			log.Println("[Parse]", "itemType not support")
			throw()
		}
	}
	cur.backToRoot()
	return cur.ptr
}

func ParseFile(filename string) (ast *node, err error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	ast = Parse(string(content))
	return
}
