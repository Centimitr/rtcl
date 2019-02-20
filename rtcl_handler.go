package rtcl

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
				args := newArgs(node.child.val)
				switch args.first {
				case "_wrapper":
					node.representation = NewContainerFromNode(node)
				case "#":
					node.representation = &Section{Name: args.second, Container: NewContainerFromNode(node)}
				case "define":
					//node.representation = &Define{Container: NewContainerFromNode(node)}
					d := &Define{Dict: make(map[string]string)}
					var k, v string
					for _, child := range astChildren(node) {
						switch child.typ {
						case "text":
							switch {
							case k == "":
								k = child.val
							case v == "":
								v = child.val
							case v != "":
								v = v + "\n" + child.val
							}
						case "blank":
							d.Dict[k] = v
							k = ""
							v = ""
						}
					}
					d.Dict[k] = v
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
				//i.representation = &Text{String: i.val}
				i.representation = "TEXT: " + i.val
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
