package rtcl

type cursor struct {
	ptr      *node
	typ      string
	onCreate func()
}

func (c *cursor) set(node *node) *cursor {
	c.ptr = node
	c.typ = node.typ
	return c
}

func newASTCursor(node *node) *cursor {
	return new(cursor).set(node)
}

func (c *cursor) clone() *cursor {
	return newASTCursor(c.ptr)
}

func (c *cursor) gotoChild() bool {
	if c.ptr.child == nil {
		return false
	}
	c.set(c.ptr.child)
	return true
}

func (c *cursor) gotoLastSibling() {
	for ; c.ptr.sibling != nil; c.set(c.ptr.sibling) {
	}
}

func (c *cursor) gotoLastChild() bool {
	if c.gotoChild() {
		c.gotoLastSibling()
		return true
	}
	return false
}

// create

func (c *cursor) signalOnCreate() {
	if c.onCreate != nil {
		c.onCreate()
	}
}

func (c *cursor) createChild(typ string) *cursor {
	if c.gotoChild() {
		c.createSibling(typ)
		return c
	}
	c.ptr.child = &node{typ: typ, parent: c.ptr}
	c.gotoChild()
	c.signalOnCreate()
	return c
}

func (c *cursor) createSibling(typ string) *cursor {
	c.gotoLastSibling()

	c.ptr.sibling = &node{typ: typ, parent: c.ptr.parent}
	c.set(c.ptr.sibling)
	c.signalOnCreate()
	return c
}

// proxy

func (c *cursor) setValue(v string) *cursor {
	c.ptr.setValue(v)
	return c
}

// back

func (c *cursor) back() *cursor {
	if c.ptr.parent == nil {
		panic("already at the highest level")
		return c
	}
	c.set(c.ptr.parent)
	return c
}

func (c *cursor) backToRoot() *cursor {
	for ; c.ptr.typ != "root"; c.set(c.ptr.parent) {
	}
	return c
}

// locate

//func preorderMatch(c *cursor, typ string, depth int) bool {
//	if c.ptr.typ == typ {
//		return true
//	}
//	if depth < 0 {
//		return false
//	}
//	cur := c.ptr
//	if cur.child != nil {
//		c.ptr = cur.child
//		if preorderMatch(c, typ, depth-1) {
//			return true
//		}
//	}
//	if cur.sibling != nil {
//		c.ptr = cur.sibling
//		if preorderMatch(c, typ, depth) {
//			return true
//		}
//	}
//	return false
//}

func (c *cursor) locate(typ string, depth int) bool {
	//return preorderMatch(c, typ, depth)
	var next []*node
	enqueue := func(node *node) {
		if node != nil {
			next = append(next, node)
		}
	}

	enqueue(c.ptr)

	for d := depth; d >= 0; d-- {
		cur := next
		next = []*node{}
		for _, node := range cur {
			switch {
			case node.child != nil && node.child.typ == typ:
				c.set(node.child)
				return true
			case node.sibling != nil && node.sibling.typ == typ:
				c.set(node.sibling)
				return true
			default:
				enqueue(node.child)
				enqueue(node.sibling)
			}
		}
	}
	return false
}

func (c *cursor) locateFromRoot(typ string, depth int) bool {
	c.backToRoot()
	return c.locate(typ, depth)
}
