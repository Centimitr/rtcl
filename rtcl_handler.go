package rtcl

import "fmt"

type Section struct {
	*Container
	Name string
}

type Define struct {
	*Container
}

type Text struct {
	String string
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
				fmt.Println(node.typ, node.val, node.child)
				if node.child == nil || node.child.typ != "block.command" {
					return
				}
				args := newArgs(node.child.val)
				switch args.first {
				case "_wrapper":
					node.representation = NewContainerFromNode(node)
				case "#":
					node.representation = &Section{Name: args.second, Container: NewContainerFromNode(node)}
				case "define":
					//node.representation = &Define{Container: NewContainerFromNode(node)}
					node.representation = "DEFINE"
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
				//i.representation = &Text{String: i.val}
				i.representation = "TEXT: " + i.val
			},
		})
}
