package main

type node struct {
	typ     string
	val     string
	parent  *node
	child   *node
	sibling *node
	ptr     *node
}

func newAST() *node {
	t := &node{typ: "root"}
	t.ptr = t
	return t
}

func (t *node) is(typ string) bool {
	return t.ptr.typ == typ
}

func (t *node) back() *node {
	if t.ptr.parent == nil {
		panic("already at the highest level")
		return t
	}
	t.ptr = t.ptr.parent
	return t
}

func (t *node) backToRoot() *node {
	for ; t.ptr != nil; t.ptr = t.ptr.parent {
	}
	return t
}

func (t *node) createChild(typ string) *node {
	if t.ptr.child != nil {
		t.ptr = t.ptr.child
		t.createSibling(typ)
		return t
	}
	t.ptr.child = &node{typ: typ, parent: t.ptr}
	t.ptr = t.ptr.child
	return t
}

func (t *node) createSibling(typ string) *node {
	for ; t.ptr.sibling != nil; t.ptr = t.ptr.sibling {
	}
	t.ptr.sibling = &node{typ: typ, parent: t.ptr.parent}
	t.ptr = t.ptr.sibling
	return t
}

func (t *node) setValue(v string) *node {
	t.ptr.val = v
	return t
}
