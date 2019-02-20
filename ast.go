package rtcl

import "github.com/rtcl/rtcl/ast"

type node struct {
	typ     string
	val     string
	parent  *node
	child   *node
	sibling *node
	//ptr            *node
	representation interface{}
}

func newAST() *node {
	return &node{typ: ast.Root}
}

func (t *node) is(typ string) bool {
	return t.typ == typ
}

func (t *node) depth() int {
	d := 0
	p := t
	for {
		if p.typ == "root" {
			break
		} else {
			p = p.parent
			d++
		}
	}
	return d
}

func (t *node) setValue(v string) *node {
	t.val = v
	return t
}

func (t *node) children() (children []*node) {
	child := t.child
	if child == nil {
		return
	}
	for {
		children = append(children, child)
		if child.sibling != nil {
			child = child.sibling
			continue
		}
		break
	}
	return
}

func astChildren(parent *node) (children []*node) {
	child := parent.child
	if child == nil {
		return
	}
	for {
		children = append(children, child)
		if child.sibling != nil {
			child = child.sibling
			continue
		}
		break
	}
	return
}
