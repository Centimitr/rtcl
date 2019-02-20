package rtcl

import (
	"fmt"
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
		Register(&Handler{
			Match: func(i *node) bool {
				return i.typ == "block.command"
			},
			Post: func(i *node) {
				i.representation = "WRAPPER TYPE: " + i.val
			},
		}).
		Register(&Handler{
			Match: func(i *node) bool {
				return i.typ == "block"
			},
			Post: func(node *node) {
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
						fmt.Println(child.typ)
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
			},
		}).
		Register(&Handler{
			Match: func(i *node) bool {
				return i.typ == "text"
			},
			Pre: func(i *node) bool {
				return false
			},
			Post: func(i *node) {
				i.representation = &Text{String: i.val}
				//i.representation = "TEXT: " + i.val
			},
		}).
		Register(&Handler{
			Match: func(i *node) bool {
				return i.typ == "paragraph"
			},
			Post: func(i *node) {
				//i.representation = "SENTENCES: " + strconv.Itoa(len(astChildren(i)))
				p := &Paragraph{}
				for _, child := range astChildren(i) {
					p.Sentences = append(p.Sentences, child.representation)
				}
				i.representation = p
			},
		})
}
