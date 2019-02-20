package rtcl

import (
	"rtcl/ast"
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
	Fragments []interface{}
	String    string
}

func (p *Paragraph) UpdateString() {
	var s string
	sep := " "
	for _, frag := range p.Fragments {
		switch v := frag.(type) {
		case *Text:
			if s != "" {
				s += sep
			}
			s += v.String
		}
	}
	p.String = s
}

type List struct {
	Options string
}

type TaskList struct {
	Tasks interface{}
}

func init() {
	DefaultHandlers.
		RegisterType(ast.Text, func(node *node, handleChildren HandleChildrenFn) {
			node.representation = &Text{String: node.val}
		}).
		RegisterType(ast.Paragraph, func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			p := &Paragraph{}
			for _, child := range astChildren(node) {
				p.Fragments = append(p.Fragments, child.representation)
			}
			p.UpdateString()
			node.representation = p
		}).

		RegisterBlockType(ast.CommandWrapper, func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			node.representation = NewContainerFromNode(node)
		}).
		RegisterBlockType("#", func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			args := NewArgsFromString(node.child.val)
			node.representation = &Section{Name: args.Second, Container: NewContainerFromNode(node)}
		}).
		RegisterBlockType("define", func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			//node.representation = &Define{Container: NewContainerFromNode(node)}
			d := &Define{Dict: make(map[string]string)}
			for _, child := range astChildren(node) {
				if child.typ == "paragraph" {
					args := Args{}
					if p, ok := child.representation.(*Paragraph); ok {
						for _, v := range p.Fragments {
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
		}).

		RegisterBlockType("[]", func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			tl := &TaskList{}
			node.representation = tl
		}).
		RegisterBlockType("-", func(node *node, handleChildren HandleChildrenFn) {
			handleChildren()
			l := &List{}
			node.representation = l
		})
}
