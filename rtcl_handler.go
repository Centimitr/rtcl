package rtcl

import (
	"strings"
)

type Section struct {
	*Container
	Name string
}

type Define struct {
	Dict map[string]string
}

type Text struct {
	String string
}

type Paragraph struct {
	Sentences []interface{}
}

func init() {
	Handlers.
		RegisterType("block.command", func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			node.representation = "WRAPPER TYPE: " + node.val
		}).
		RegisterType("block", func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			if node.child == nil || node.child.typ != "block.command" {
				return
			}
			args := NewArgsFromString(node.child.val)
			switch args.First {
			case "_wrapper":
				node.representation = NewContainerFromNode(node)
			case "#":
				node.representation = &Section{Name: args.Second, Container: NewContainerFromNode(node)}
			case "define":
				//node.representation = &Define{Container: NewContainerFromNode(node)}
				d := &Define{Dict: make(map[string]string)}
				for _, child := range astChildren(node) {
					if child.typ == "paragraph" {
						args := Args{}
						if p, ok := child.representation.(*Paragraph); ok {
							for _, v := range p.Sentences {
								if text, ok := v.(*Text); ok {
									args.Append(text.String)
								} else {
									panic("define: paragraphs in define only accept 'text' type")
								}
							}
						}
						if args.First != "" {
							d.Dict[args.First] = strings.Join(args.Slice[1:], "\n")
						}
					}
				}
				node.representation = d
			}
		}).
		RegisterType("text", func(node *node, handleChildren HandleChildrenFn) {
			//i.representation = "TEXT: " + i.val
			node.representation = &Text{String: node.val}
		}).
		RegisterType("paragraph", func(node *node, handleChildren HandleChildrenFn) {
			//i.representation = "SENTENCES: " + strconv.Itoa(len(astChildren(i)))
			handleChildren()
			p := &Paragraph{}
			for _, child := range astChildren(node) {
				p.Sentences = append(p.Sentences, child.representation)
			}
			node.representation = p
		})
}
