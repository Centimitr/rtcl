package main

type node struct {
	typ            string
	val            string
	parent         *node
	child          *node
	sibling        *node
	ptr            *node
	representation interface{}
}

func newAST() *node {
	t := &node{typ: "root"}
	t.ptr = t
	return t
}

func (t *node) is(typ string) bool {
	return t.ptr.typ == typ
}

func (t *node) depth() int {
	d := 0
	p := t.ptr
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

func (t *node) back() *node {
	if t.ptr.parent == nil {
		panic("already at the highest level")
		return t
	}
	t.ptr = t.ptr.parent
	return t
}

func (t *node) backToRoot() *node {
	for ; t.ptr.typ != "root"; t.ptr = t.ptr.parent {
	}
	return t
}

func (t *node) onCreate() {
	//fmt.Println(strings.Repeat("    ", t.depth()-1) + t.ptr.typ)
}

func (t *node) createChild(typ string) *node {
	if t.ptr.child != nil {
		t.ptr = t.ptr.child
		t.createSibling(typ)
		return t
	}
	t.ptr.child = &node{typ: typ, parent: t.ptr}
	t.ptr = t.ptr.child
	t.onCreate()
	return t
}

func (t *node) createSibling(typ string) *node {
	for ; t.ptr.sibling != nil; t.ptr = t.ptr.sibling {
	}
	t.ptr.sibling = &node{typ: typ, parent: t.ptr.parent}
	t.ptr = t.ptr.sibling
	t.onCreate()
	return t
}

func (t *node) setValue(v string) *node {
	t.ptr.val = v
	return t
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

func preorderMatch(ast *node, typ string, depth int) bool {
	if ast.ptr.typ == typ {
		return true
	}
	if depth < 0 {
		return false
	}
	cur := ast.ptr
	if cur.child != nil {
		ast.ptr = cur.child
		if preorderMatch(ast, typ, depth-1) {
			return true
		}
	}
	if cur.sibling != nil {
		ast.ptr = cur.sibling
		if preorderMatch(ast, typ, depth) {
			return true
		}
	}
	return false
}

func (t *node) locate(typ string, depth int) bool {
	return preorderMatch(t, typ, depth)
}

func (t *node) locateFromRoot(typ string, depth int) bool {
	t.backToRoot()
	return t.locate(typ, depth)
}
